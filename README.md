# golang-sonzai-bank

4 mức độ cô lâp data trong SQL

- Read Uncommitted : 
  - khi 2 transaction cùng select đến 1 hàng dữ liệu, sau đó transaction 1 update hàng dữ liệu đó 
  nhưng chưa commit. thì transacion 2 select thì vẫn nhận được dẽ liệu mà transaction thay đổi dù chưa commit
  -> transaction 2 select ra dữ liệu bẩn (dirty read)
- Read commited
- Repeatable Read:
  - khi transaction đã bắt đầu thì các dữ liệu select sẽ không đổi. kể cả khi có 1 transaction khác 
  đã thay đổi và commit dữ liệu mới. 
  nếu transaction 2 thay đổi data ( balance -10 ) thì nó vẫn lấy được balance đúng và trừ đi 10 
- serializable : 