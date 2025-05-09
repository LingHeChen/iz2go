package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

// DeclarationType 表示声明类型
type DeclarationType int

const (
	VarDeclaration DeclarationType = iota
	ConstDeclaration
	TypeDeclaration
	FuncDeclaration
	InterfaceDeclaration
	StructDeclaration
	MethodDeclaration // 新增：方法声明
)

// CheckDeclaration 检查文件中是否存在指定的声明
// filePath: 文件路径
// name: 要检查的声明名称
// declType: 声明类型
// receiverType: 方法接收者类型
// 返回: 是否存在该声明
func CheckDeclaration(filePath string, name string, declType DeclarationType, receiverType string) (bool, error) {
	// 创建文件集
	fset := token.NewFileSet()

	// 解析文件
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return false, fmt.Errorf("解析文件失败: %v", err)
	}

	// 检查声明是否存在
	var found bool
	ast.Inspect(f, func(n ast.Node) bool {
		switch declType {
		case VarDeclaration:
			if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.VAR {
				for _, spec := range decl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for _, ident := range valueSpec.Names {
							if ident.Name == name {
								found = true
								return false
							}
						}
					}
				}
			}

		case ConstDeclaration:
			if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.CONST {
				for _, spec := range decl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for _, ident := range valueSpec.Names {
							if ident.Name == name {
								found = true
								return false
							}
						}
					}
				}
			}

		case TypeDeclaration:
			if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
				for _, spec := range decl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if typeSpec.Name.Name == name {
							found = true
							return false
						}
					}
				}
			}

		case FuncDeclaration:
			if funcDecl, ok := n.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name == name {
					found = true
					return false
				}
			}

		case InterfaceDeclaration:
			if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
				for _, spec := range decl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if typeSpec.Name.Name == name {
							if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
								found = true
								return false
							}
						}
					}
				}
			}

		case StructDeclaration:
			if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
				for _, spec := range decl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if typeSpec.Name.Name == name {
							if _, ok := typeSpec.Type.(*ast.StructType); ok {
								found = true
								return false
							}
						}
					}
				}
			}

		case MethodDeclaration:
			if funcDecl, ok := n.(*ast.FuncDecl); ok && funcDecl.Recv != nil {
				// 检查接收者类型
				for _, field := range funcDecl.Recv.List {
					if typeExpr, ok := field.Type.(*ast.StarExpr); ok {
						// 处理指针接收者
						if ident, ok := typeExpr.X.(*ast.Ident); ok && ident.Name == receiverType {
							if funcDecl.Name.Name == name {
								found = true
								return false
							}
						}
					} else if ident, ok := field.Type.(*ast.Ident); ok {
						// 处理值接收者
						if ident.Name == receiverType {
							if funcDecl.Name.Name == name {
								found = true
								return false
							}
						}
					}
				}
			}
		}
		return true
	})

	return found, nil
}

// CheckDeclarationInDir 检查目录中所有 Go 文件是否存在指定的声明
// dirPath: 目录路径
// name: 要检查的声明名称
// declType: 声明类型
// 返回: 是否存在该声明
func CheckDeclarationInDir(dirPath string, name string, declType DeclarationType) (bool, error) {
	// 获取目录下所有 Go 文件
	files, err := filepath.Glob(filepath.Join(dirPath, "*.go"))
	if err != nil {
		return false, fmt.Errorf("获取文件列表失败: %v", err)
	}

	// 检查每个文件
	for _, file := range files {
		found, err := CheckDeclaration(file, name, declType, "")
		if err != nil {
			return false, err
		}
		if found {
			return true, nil
		}
	}

	return false, nil
}

// CheckMethod 检查单个文件中的方法
func CheckMethod(filePath string, receiverType string, methodName string) (bool, error) {
	return CheckDeclaration(filePath, methodName, MethodDeclaration, receiverType)
}

// CheckMethodInDir 检查目录中所有 Go 文件是否存在指定的方法
// dirPath: 目录路径
// receiverType: 方法接收者类型
// methodName: 方法名称
// 返回: 是否存在该方法
func CheckMethodInDir(dirPath string, receiverType string, methodName string) (bool, error) {
	// 获取目录下所有 Go 文件
	files, err := filepath.Glob(filepath.Join(dirPath, "*.go"))
	if err != nil {
		return false, fmt.Errorf("获取文件列表失败: %v", err)
	}

	// 检查每个文件
	for _, file := range files {
		found, err := CheckMethod(file, receiverType, methodName)
		if err != nil {
			return false, err
		}
		if found {
			return true, nil
		}
	}

	return false, nil
}

func CheckConst(filePath string, name string) (bool, error) {
	return CheckDeclaration(filePath, name, ConstDeclaration, "")
}
