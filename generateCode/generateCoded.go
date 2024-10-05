package generateCode

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Field struct {
	Name     string
	Type     string
	Nullable bool
	GormTag  string // GORM 标签
	Validate string // 验证标签
	JsonTag  string // JSON 标签
}

type Table struct {
	Name   string
	Fields []Field
}

// 生成结构体数据
func generateModel(table Table) string {
	model := fmt.Sprintf("package model\n\nimport \"gorm.io/gorm\"\n\ntype %s struct {\n", capitalize(table.Name))
	model += "    gorm.Model\n" // 包含默认字段
	for _, field := range table.Fields {
		// 根据 Nullable 生成 GORM 和 Validate 标签
		gormTag := field.GormTag
		if !field.Nullable {
			gormTag += ";not null"
		}

		validateTag := ""
		if !field.Nullable {
			validateTag = "required"
		}
		tags := fmt.Sprintf("`gorm:\"%s\" validate:\"%s\" json:\"%s\"`", gormTag, validateTag, field.JsonTag)
		model += fmt.Sprintf("    %s %s %s\n", capitalize(field.Name), mapSQLTypeToGoType(field.Type), tags)
	}
	model += "}\n"
	return model
}

// 映射数据类型
func mapSQLTypeToGoType(sqlType string) string {
	switch sqlType {
	case "int":
		return "int"
	case "varchar":
		return "string"
	case "text":
		return "string"
	case "datetime":
		return "time.Time"
	case "float":
		return "float64"
	// 添加更多 SQL 类型映射
	default:
		return "interface{}"
	}
}

// 将第一个字母转为大写
func capitalize(str string) string {
	if len(str) == 0 {
		return str
	}
	return string(str[0]-32) + str[1:] // 将首字母大写
}

// 将字段名称转换为驼峰式
func ToCamelCase(snake string) string {
	parts := strings.Split(snake, "_")
	for i := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(parts[i]) // 第一个单词小写
		} else {
			parts[i] = strings.Title(parts[i]) // 其余单词首字母大写
		}
	}
	return strings.Join(parts, "")
}

// 生成简单crud文件
func generateCRUD(table Table) string {
	modelName := capitalize(table.Name)
	return fmt.Sprintf(`package repository

import (
    "gorm.io/gorm"
)

func Create%s(db *gorm.DB, record *%s) error {
    return db.Create(record).Error
}

func Get%s(db *gorm.DB, id uint) (*%s, error) {
    var record %s
    if err := db.First(&record, id).Error; err != nil {
        return nil, err
    }
    return &record, nil
}

func Update%s(db *gorm.DB, record *%s) error {
    return db.Save(record).Error
}

func Delete%s(db *gorm.DB, id uint) error {
    return db.Delete(&%s{}, id).Error
}
`, modelName, modelName, modelName, modelName, modelName, modelName, modelName, modelName, modelName)
}

// 保存到文件
func saveToFile(filePath string, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}

func Generate(table Table) error {

	//// 定义表及字段
	//table := Table{
	//	Name: "user",
	//	Fields: []Field{
	//		{Name: "account_id", Type: "int", Nullable: false, GormTag: "comment:'账号ID，关联account表'", Validate: "required", JsonTag: "account_id,omitempty"},
	//		{Name: "name", Type: "varchar", Nullable: false, GormTag: "comment:'用户名'", Validate: "required", JsonTag: "name,omitempty"},
	//		{Name: "email", Type: "varchar", Nullable: true, GormTag: "comment:'邮箱'", Validate: "email", JsonTag: "email,omitempty"},
	//	},
	//}

	// 生成结构体代码
	modelCode := generateModel(table)
	//fmt.Println(modelCode)

	// 生成 CRUD 代码
	crudCode := generateCRUD(table)
	//fmt.Println(crudCode)
	// 创建目录
	modelDir := "./generate/model"
	repoDir := "./generate/repository"

	os.MkdirAll(modelDir, os.ModePerm)
	os.MkdirAll(repoDir, os.ModePerm)

	// 保存结构体到文件
	modelFilePath := filepath.Join(modelDir, fmt.Sprintf("%s.go", table.Name))
	if err := saveToFile(modelFilePath, modelCode); err != nil {
		//return //fmt.Errorf("Error saving model file: %v\n", err)
		//} else {
		return fmt.Errorf("Model saved to %s\n", modelFilePath)
	}

	// 保存 CRUD 操作到文件
	crudFilePath := filepath.Join(repoDir, fmt.Sprintf("%s_repository.go", table.Name))
	if err := saveToFile(crudFilePath, crudCode); err != nil {
		//return fmt.Errorf("Error saving CRUD file: %v\n", err)
		//} else {
		return fmt.Errorf("CRUD saved to %s\n", crudFilePath)
	}
	return nil
}
