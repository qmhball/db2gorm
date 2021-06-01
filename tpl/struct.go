package tpl

const(
    StructTpl = `
package {{.PackageName}}

type {{.StructName}} struct{
{{- range $i, $v := .ColumnsInfo}}
    {{$v.Field}} {{$v.Type}} {{$v.Default}}
{{- end}}
}
`
)