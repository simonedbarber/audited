# Audited

Audited is used to record the last User who created and/or updated your [GORM](https://gorm.io/gorm) model. It does so using a `CreatedBy` and `UpdatedBy` field. While Audited can be used alone (with [GORM](https://gorm.io/gorm)), it integrates nicely with [QOR](https://github.com/simonedbarber/qor) to log and display this extra information.

[![GoDoc](https://godoc.org/github.com/simonedbarber/audited?status.svg)](https://godoc.org/github.com/simonedbarber/audited)

### Register GORM Callbacks

Audited utilizes [GORM](https://gorm.io/gorm) callbacks to log data, so you will need to register callbacks first:

```go
import (
  "gorm.io/gorm"
  "github.com/simonedbarber/audited"
)

db, err := gorm.Open("sqlite3", "demo_db")
audited.RegisterCallbacks(db)
```

### Making a Model Auditable

Embed `audited.AuditedModel` into your model as an anonymous field to make the model auditable:

```go
type Product struct {
  gorm.Model
  Name string
  audited.AuditedModel
}
```

### Usage

```go
import "github.com/simonedbarber/audited"
import "gorm.io/gorm"

func main() {
  var db, err = gorm.Open("sqlite3", "demo_db")
  var currentUser = User{ID: 100}
  var product Product

  // Create will set product's `CreatedBy`, `UpdatedBy` to `currentUser`'s primary key if `audited:current_user` is a valid model
  db.Set("audited:current_user", currentUser).Create(&product)
  // product.CreatedBy => 100
  // product.UpdatedBy => 100

  // If it is not a valid model, then will set `CreatedBy`, `UpdatedBy` to the passed value
  db.Set("audited:current_user", "admin").Create(&product)
  // product.CreatedBy => "admin"
  // product.UpdatedBy => "admin"

  // When updating a record, it will update the `UpdatedBy` to `audited:current_user`'s value
  db.Set("audited:current_user", "dev").Model(&product).Update("Code", "L1212")
  // product.UpdatedBy => "dev"
}
```

[QOR Demo:  http://demo.getqor.com/admin](http://demo.getqor.com/admin)

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).
