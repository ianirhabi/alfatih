# astaxie/beego/orm with extra steroid


## Installation
```
    go get git.qasico.com/cuxs/orm
```

## Set Up Database
ORM supports three popular databases. Here are the tested drivers, you need to import them:
```
import (
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"
    _ "github.com/mattn/go-sqlite3"
)
```
## Func in Cuxs ORM

### Ormer interface
define the orm interface

- **Read(md interface{}, cols ...string) error**<br />
  read data to model

- **ReadForUpdate(md interface{}, cols ...string) error**<br />
  Like Read(), but with "FOR UPDATE" clause, useful in transaction.
  Some databases are not support this feature.

- **ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error)**<br />
  Try to read a row from the database, or insert one if it doesn't exist

- **Insert(interface{}) (int64, error)**<br />
  insert model data to database.
  user must a pointer and Insert will set user's pk field

- **InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error)**<br />
  mysql:InsertOrUpdate(model) or InsertOrUpdate(model,"colu=colu+value"). <br />
  if column type is integer : can use(+-*/), string : convert(colu,"value")

- **InsertMulti(bulk int, mds interface{}) (int64, error)**<br />
  insert some models to database

- **Update(md interface{}, cols ...string) (int64, error)**<br />
  update model to database. cols set the columns those want to update.<br />
  find model by Id(pk) field and update columns specified by fields, if cols is null then update all columns

- **Delete(md interface{}, cols ...string) (int64, error)**<br />
  delete model in database

- **LoadRelated(md interface{}, name string, args ...interface{}) (int64, error)**<br />
  load related models to md model. args are limit, offset int and order string.<br />
  make sure the relation is defined in model struct tags.<br />
  example:<br />
  ***orm.LoadRelated(post,"Tags")***<br />
  ***for _,tag := range post.Tags{...}***<br />

- **QueryM2M(md interface{}, name string) QueryM2Mer**<br />
  create a models to models queryer

- **QueryTable(ptrStructOrTableName interface{}) QuerySeter**<br />
  return a QuerySeter for table operations. table name can be string or struct.<br />
  e.g. QueryTable("user"), QueryTable(&user{}) or QueryTable((*User)(nil))

- **Using(name string) error**<br />
  switch to another registered database driver by given name.

- **Begin() error**<br />
  begin transaction

- **Commit() error**<br />
  commit transaction

- **Rollback() error**<br />
  rollback transaction

- **Raw(query string, args ...interface{}) RawSeter**<br />
  return a raw query seter for raw sql string.

- **Driver() Driver**<br />

### Inserter interface
Inserter insert prepared statement
- **Insert(interface{}) (int64, error)**<br />
- **Close() error**

### QuerySeter interface
QuerySeter query seter

- **Filter(string, ...interface{}) QuerySeter**<br />
  add condition expression to QuerySeter.

- **Exclude(string, ...interface{}) QuerySeter**<br />
  add NOT condition to querySeter. have the same usage as Filter

- **SetCond(*Condition) QuerySeter**<br />
  set condition to QuerySeter. sql's where condition

- **Limit(limit interface{}, args ...interface{}) QuerySeter**<br />
  add LIMIT value. args[0] means offset, e.g. LIMIT num,offset.<br />
  if Limit <= 0 then Limit will be set to default limit ,eg 1000<br />
  if QuerySeter doesn't call Limit, the sql's Limit will be set to default limit, eg 1000<br />
   for example:<br />
	***qs.Limit(10, 2)***<br />
	***// sql-> limit 10 offset 2***

- **Offset(offset interface{}) QuerySeter**<br />
  add OFFSET value

- **GroupBy(exprs ...string) QuerySeter**<br />
  add GROUP BY expression

- **OrderBy(exprs ...string) QuerySeter**<br />
  add ORDER expression.

- **RelatedSel(params ...interface{}) QuerySeter**<br />
  set relation model to query together.it will query relation models and assign to parent model.<br />
  for example:<br />
	***// will load all related fields use left join.***<br />
	***qs.RelatedSel().One(&user)***<br />
	***// will  load related field only profile***<br />
	***qs.RelatedSel("profile").One(&user)***<br />
	***user.Profile.Age = 32***


- **Distinct() QuerySeter**<br />
  Set Distinct

- **Count() (int64, error)**<br />
  return QuerySeter execution result number

- **Exist() bool**<br />
  check result empty or not after QuerySeter executed

- **Update(values Params) (int64, error)**<br />
  execute update with parameters

- **Delete() (int64, error)**<br />
  delete from table

- **PrepareInsert() (Inserter, error)**<br />
  return a insert queryer. it can be used in times.

- **All(container interface{}, cols ...string) (int64, error)**<br />
  query all data and map to containers. cols means the columns when querying.

- **One(container interface{}, cols ...string) error**<br />
  query one row data and map to containers. cols means the columns when querying.<br />
  for example:<br />
	***var user User***<br />
	***qs.One(&user) //user.UserName == "slene"***

- **Values(results *[]Params, exprs ...string) (int64, error)**<br />
  it converts data to []map[column]value.

- **ValuesList(results *[]ParamsList, exprs ...string) (int64, error)**<br />
  query all data and map to [][]interface. it converts data to [][column_index]value

- **ValuesFlat(result *ParamsList, expr string) (int64, error)**<br />
  query all data and map to []interface.it's designed for one column record set, auto change to []value, not [][column]value.

- **RowsToMap(result *Params, keyCol, valueCol string) (int64, error)**<br />
  query all rows into map[string]interface with specify key and value column name.

- **RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error)**<br />
  query all rows into struct with specify key and value column name.


### QueryM2Mer interface
QueryM2Mer model to model query struct. all operations are on the m2m table only and will not affect the origin model table

- **Add(...interface{}) (int64, error)**<br />
  add models to origin models when creating queryM2M. insert one or more rows to m2m table<br />
  make sure the relation is defined in post model struct tag.

- **Remove(...interface{}) (int64, error)**<br />
  remove models following the origin model relationship. Only delete rows from m2m table

- **Exist(interface{}) bool**<br />
  check model is existed in relationship of origin model

- **Clear() (int64, error)**<br />
  clean all models in related of origin model

- **Count() (int64, error)**<br />
  count all related models of origin model


### RawPreparer interface
RawPreparer raw query statement
- **Exec(...interface{}) (sql.Result, error)**
- **Close() error**

### RawSeter interface
RawSeter raw query seter, create From Ormer.Raw<br />
for example:<br />
***sql := fmt.Sprintf("SELECT %sid%s,%sname%s FROM %suser%s WHERE id = ?",Q,Q,Q,Q,Q,Q)***<br />
***rs := Ormer.Raw(sql, 1)***

- **Exec() (sql.Result, error)**<br />
  execute sql and get result

- **QueryRow(containers ...interface{}) error**<br />
  query data and map to container

- **QueryRows(containers ...interface{}) (int64, error)**<br />
  query data rows and map to container

- **SetArgs(...interface{}) RawSeter**<br />
  set args for every query

- **Values(container *[]Params, cols ...string) (int64, error)**<br />
  query data to []map[string]interface. see QuerySeter's Values

- **ValuesList(container *[]ParamsList, cols ...string) (int64, error)**<br />
  query data to [][]interface. see QuerySeter's ValuesList

- **ValuesFlat(container *ParamsList, cols ...string) (int64, error)**<br />
  query data to []interface. see QuerySeter's ValuesFlat

- **RowsToMap(result *Params, keyCol, valueCol string) (int64, error)**<br />
  query all rows into map[string]interface with specify key and value column name.

- **RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error)**<br />
  query all rows into struct with specify key and value column name.

- **Prepare() (RawPreparer, error)**<br />
  return prepared raw statement for used in times.

### Condition ORM
Condition struct
```go

type condValue struct {
	exprs  []string
	args   []interface{}
	cond   *Condition
	isOr   bool
	isNot  bool
	isCond bool
}
// work for WHERE conditions.
type Condition struct {
	params []condValue
}
```
Condition func:
- **func NewCondition() *Condition**<br />
  NewCondition return new condition struct
- **func (c Condition) And(expr string, args ...interface{}) *Condition**<br />
  And add expression to condition
- **func (c Condition) AndNot(expr string, args ...interface{}) *Condition**<br />
  AndNot add NOT expression to condition
- **func (c *Condition) AndCond(cond *Condition) *Condition**<br />
  AndCond combine a condition to current condition
- **func (c Condition) Or(expr string, args ...interface{}) *Condition**<br />
  add OR expression to condition
- **func (c Condition) OrNot(expr string, args ...interface{}) *Condition**<br />
  add OR NOT expression to condition
- **func (c *Condition) OrCond(cond *Condition) *Condition**<br />
  combine a OR condition to current condition
- **func (c *Condition) IsEmpty() bool**<br />
  check the condition arguments are empty or not.
- **func (c Condition) clone() *Condition** <br />
  clone a condition


## Query Builder
ORM is more for simple CRUD operations, whereas QueryBuilder is for complex queries with subqueries and multi-joins.<br />
The list for QueryBuilder objects are below:
```
type QueryBuilder interface {
	Select(fields ...string) QueryBuilder // Select will join the fields
	From(tables ...string) QueryBuilder // From join the tables
	InnerJoin(table string) QueryBuilder // InnerJoin INNER JOIN the table
	LeftJoin(table string) QueryBuilder // LeftJoin LEFT JOIN the table
	RightJoin(table string) QueryBuilder // RightJoin RIGHT JOIN the table
	On(cond string) QueryBuilder // On join with on cond
	Where(cond string) QueryBuilder // Where join the Where cond
	And(cond string) QueryBuilder // And join the and cond
	Or(cond string) QueryBuilder // Or join the or cond
	In(vals ...string) QueryBuilder // In join the IN (vals)
	OrderBy(fields ...string) QueryBuilder // OrderBy join the Order by fields
	Asc() QueryBuilder // Asc join the asc
	Desc() QueryBuilder // Desc join the desc
	Limit(limit int) QueryBuilder // Limit join the limit num
	Offset(offset int) QueryBuilder // Offset join the offset num
	GroupBy(fields ...string) QueryBuilder // GroupBy join the Group by fields
	Having(cond string) QueryBuilder // Having join the Having cond
	Subquery(sub string, alias string) string // Subquery join the sub as alias
	String() string // String join all Tokens
}
//Func NewQueryBuilder return the QueryBuilder
// only some driver are available
func NewQueryBuilder(driver string) (qb QueryBuilder, err error)
```

## Set parameters in ORM
#### Relation
Use `;` as the separator of multiple settings. Use `,` as the separator if a setting has multiple values.
```
orm:"null;rel(fk)"
```
#### Ignore Field
Use `-` to ignore field in the struct.
```
type User struct {

    AnyField string `orm:"-"`

}
```
#### Auto
When Field type is int, int32, int64, uint, uint32 or uint64, you can set it as auto increment.<br />
If there is no primary key in the model definition, the field `Id` with one of the types above will be considered as auto increment key<br />
Because of the design of go.

#### pk
Set as primary key. Used for using other type field as primary key.

#### null
Fields are `NOT NULL` by default. Set null to `ALLOW NULL`.
```go
Name string `orm:"null"`
```

#### index
Add index for one field

#### unique
Add unique key for one field
```go
Name string `orm:"unique"`
```

#### column
Set column name in db table for field.
```go
Name string `orm:"column(user_name)"`
```

#### size
Default value for string field is varchar(255). It will use varchar(size) after setting.
```go
Title string `orm:"size(60)"`
```

#### digits / decimals
Set precision for float32 or float64.
```go
Money float64 `orm:"digits(12);decimals(4)"`
```
Total 12 digits, 4 digits after point. For example: `12345678.1234`

#### auto_now / auto_now_add
```go
Created time.Time `orm:"auto_now_add;type(datetime)"`
Updated time.Time `orm:"auto_now;type(datetime)"`
```
* auto_now: every save will update time.
* auto_now_add: set time at the first save
This setting won't affect massive `update`.

#### type
If set type as date, the field's db type is date.
```go
Created time.Time `orm:"auto_now_add;type(date)"`
```
If set type as datetime, the field's db type is datetime.
```go
Created time.Time `orm:"auto_now_add;type(datetime)"`
```

#### default
Set default value for field with the same type. (Only support default value of cascade deleting.)
```go
type User struct {
    ...
    Status int `orm:"default(1)"`
    ...
}
```

## Relationship Table

#### rel / reverse
**RelOneToOne**:
```go
type User struct {
    ...
    Profile *Profile `orm:"null;rel(one);on_delete(set_null)"`
    ...
}
```

**RelForeignKey**:
```go
type Post struct {
    ...
    User *User `orm:"rel(fk)"` // RelForeignKey relation
    ...
}
```

**RelManyToMany**:
```go
type Post struct {
    ...
    Tags []*Tag `orm:"rel(m2m)"` // ManyToMany relation
    ...
}
```

The reverse relationship <br />
**RelReverseOne**:
```go
type Profile struct {
    ...
    User *User `orm:"reverse(one)"`
    ...
}
```

**RelReverseMany**:
```go
type Tag struct {
    ...
    Posts []*Post `orm:"reverse(many)"` // reverse relationship of fk
    ...
}
```

### example of set parameter
```go
type test struct {
	ID            int64      `orm:"column(id);auto"`
	Users         *User      `orm:"column(User_id);rel(fk)"`
	City          string     `orm:"column(city);size(50)"`
	List         []*Name     `orm:"reverse(many)"`
}
```

It's field type depends on related primary key.

* RelForeignKey
* RelOneToOne
* RelManyToMany
* RelReverseOne
* RelReverseMany

## Example ORM
### Insert,Update and Delete Usage
```go
func init() {
	// RegisterModel register models
	orm.RegisterModel(new(tes))
}

type tes struct {
	ID            int64      `orm:"column(id);auto"`
	User          string     `orm:"column(user)"`
	Password      string     `orm:"column(pass)"`
}
// It will updating if this struct has valid Id
// if not, will inserting a new row.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (t *tes) Save(field ...string) (err error) {
	o := orm.NewOrm()
	if t.ID > 0 {
		// update table in database
		_, err = o.Update(t, field...)
	} else {
		// insert table in database
		t.ID, err = o.Insert(t)
	}
	return
}
// Delete permanently deleting tes data
func (t *tes) Delete() error {
	o := orm.NewOrm()
	if t.ID > 0 {
		_, err := o.Delete(t)
		return err
	}
	return orm.ErrNoRows
}
```
### Get All Data and Single Data
```go
// Gettes find a single data tes using field and value condition.
func Gettes(field string, values ...interface{}) (*model.tes, error) {
	// The new built-in function allocates memory. The first argument is a type,
	// not a value, and the value returned is a pointer to a newly
	// allocated zero value of that type(model.tes).
	m := new(model.tes)
	// create new orm and return a QuerySeter for table operations.
	o := orm.NewOrm().QueryTable(m)
	// using Filter to add condition expression to QuerySeter.
	// set relation model to query together.
	// add LIMIT value to 1 (only get 1 row)
 	// query one row data and map to containers.
	if err := o.Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	// return m container and error nil
	return m, nil
}
// Gettess  get all data tes that matched with query request parameters.
func Gettess(rq *orm.RequestQuery) (*[]model.tes, int64, error) {
	var m []model.tes
	var count int64
	var err error

	o := orm.NewOrm()
	qs := o.QueryTable(new(model.tes)).SetCond(rq.GetCondition())

	if len(rq.Embeds) > 0 {
		qs = qs.RelatedSel(rq.GetJoin())
	}

	if count, err = qs.Count(); err != nil || count == 0 {
		return nil, count, err
	}
        // add ORDER expression
	qs = qs.OrderBy(rq.OrderBy...).Limit(rq.Limit, rq.Offset)
	// query all data and map to containers.
	if _, err = qs.All(&m, rq.Fields...); err == nil {
		return &m, count, nil
	}
	return nil, count, err
}
```

***Get All Related Data***<br />
Gettesrelated get all data including fk data and return all data, amount of data and error
```go
func Gettesrelated(rq *orm.RequestQuery) (*[]model.tes, int64, error) {
	var m []model.tes
	var count int64
	var err error

	o := orm.NewOrm()
	qs := o.QueryTable(new(model.tes)).SetCond(rq.GetCondition())

	if len(rq.Embeds) > 0 {
		qs = qs.RelatedSel(rq.GetJoin())
	}

	if count, err = qs.Count(); err != nil || count == 0 {
		return nil, count, err
	}

	qs = qs.OrderBy(rq.OrderBy...).Limit(rq.Limit, rq.Offset)
	if _, err = qs.All(&m, rq.Fields...); err == nil {
		var mx []model.tes
		for _, tes := range m {
			o.LoadRelated(&tes, "Tags")
			mx = append(mx, tes)

		}
		return &mx, count, nil
	}

	return nil, count, err
}
```