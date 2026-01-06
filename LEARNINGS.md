# sqlc tool
1. With `:one`, if no row is updated, `err` is thrown. No error is thrown
   though for `:exec` or `:execrows` but you can check the rows affected
   returned by `:execrows` to determine this. No way exists for `:exec`.
