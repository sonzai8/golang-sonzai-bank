package gapi

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/sonzai8/golang-sonzai-bank/db/sqlc"
	"github.com/sonzai8/golang-sonzai-bank/pb"
	"github.com/sonzai8/golang-sonzai-bank/utils"
	"github.com/sonzai8/golang-sonzai-bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	authPayload, err := server.AuthorizeUser(ctx)
	if err != nil {
		return nil, unAuthenticationError(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other users info")
	}

	arg := db.UpdateUserParams{
		HashedPassword: pgtype.Text{},
		FullName: pgtype.Text{
			String: req.GetFullName(),
			Valid:  req.GetFullName() != "",
		},
		Email: pgtype.Text{
			String: req.GetEmail(),
			Valid:  req.GetEmail() != "",
		},
		Username: req.GetUsername(),
	}

	if req.GetPassword() != "" {
		hashedPassword, err := utils.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to hash password: %v", err)
		}
		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}
	}
	account, err := server.store.UpdateUser(ctx, arg)
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

	rspUSer := &pb.UpdateUserResponse{
		User: convertUser(account),
	}

	return rspUSer, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if req.GetFullName() != "" {
		if err := val.ValidateFullname(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}
	if req.GetEmail() != "" {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}
	if req.GetPassword() != "" {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	return violations
}
