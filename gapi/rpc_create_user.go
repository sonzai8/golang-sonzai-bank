package gapi

import (
	"context"
	"errors"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/sonzai8/golang-sonzai-bank/db/sqlc"
	"github.com/sonzai8/golang-sonzai-bank/pb"
	"github.com/sonzai8/golang-sonzai-bank/utils"
	"github.com/sonzai8/golang-sonzai-bank/val"
	"github.com/sonzai8/golang-sonzai-bank/worker"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password: %v", err)

	}
	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			// TODO: user db transaction
			// Send Verify email here
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: req.GetUsername(),
			}

			taskOpts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, taskOpts...)
		},
	}
	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			//fmt.Println("check error code: ", pgErr.Code)
			switch pgErr.Code {
			case "23505":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists")
			}
		}
		return nil, status.Errorf(codes.Internal, "Failed to create user: %v", err)

	}

	rspUSer := &pb.CreateUserResponse{
		User: convertUser(txResult.User),
	}

	return rspUSer, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := val.ValidateFullname(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	return violations
}
