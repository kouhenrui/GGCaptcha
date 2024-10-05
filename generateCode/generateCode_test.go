package generateCode

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	// 定义表及字段

	// 定义表及字段
	table := Table{
		Name: "account",
		Fields: []Field{
			{Name: ToCamelCase("phone"), Type: "int", Nullable: false, GormTag: "comment:'phone'", Validate: "required", JsonTag: "phone,omitempty"},
			{Name: ToCamelCase("name"), Type: "varchar", Nullable: false, GormTag: "comment:'用户名'", Validate: "required", JsonTag: "name,omitempty"},
			{Name: ToCamelCase("email"), Type: "varchar", Nullable: true, GormTag: "comment:'邮箱'", Validate: "email", JsonTag: "email,omitempty"},
			{Name: ToCamelCase("password"), Type: "varchar", Nullable: false, GormTag: "comment:'hash密码'", JsonTag: "password"},
			{Name: ToCamelCase("salt"), Type: "varchar", Nullable: false, GormTag: "comment:'加密盐'", JsonTag: "salt"},
		},
	}
	err := Generate(table)
	if err != nil {
		t.Fatalf("生成文件错误%s", err)
	}
	t.Log("生成成功")
	//table := Table{
	//	Name: "user",
	//	Fields: []Field{
	//		{Name: "account_id", Type: "int", Nullable: false, GormTag: "comment:'账号ID，关联account表'", Validate: "required", JsonTag: "account_id,omitempty"},
	//		{Name: "name", Type: "varchar", Nullable: false, GormTag: "comment:'用户名'", Validate: "required", JsonTag: "name,omitempty"},
	//		{Name: "email", Type: "varchar", Nullable: true, GormTag: "comment:'邮箱'", Validate: "email", JsonTag: "email,omitempty"},
	//	},
	//}
	//
	//// 生成结构体和 CRUD 代码
	//modelCode := generateModel(table)
	//t.Log(modelCode)
	////fmt.Println(modelCode)
	//t.Log("--------------------------------------------")
	//crudCode := generateCRUD(table)
	//t.Log(crudCode)
	//fmt.Println(crudCode)
}
