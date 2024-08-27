package main

import (
	"fmt"
	"log"

	"github.com/forhsd/postgresql-parser/pkg/sql/parser"
	"github.com/forhsd/postgresql-parser/pkg/sql/sem/tree"
	"github.com/forhsd/postgresql-parser/pkg/walk"
)

func main() {
	sql := `select marr,1::text id
			from (select marr_stat_cd AS marr, label AS l
				  from public.root_loan_mock_v4 aa
				  left join public.user cc on aa.id = cc.id
				  order by root_loan_mock_v4.age desc, l desc
				  limit 5) as v4
			LIMIT 1;`
	w := &walk.AstWalker{
		Fn: func(ctx interface{}, node interface{}) (stop bool) {
			// log.Printf("node type %T %+v", node, node)

			te, ok := node.(*tree.AliasedTableExpr)
			if ok {
				tn, exists := te.Expr.(*tree.TableName)
				if exists {
					if tn.TableNamePrefix.ExplicitSchema {
						fmt.Printf("实体表: %v.%v\n", tn.TableNamePrefix.SchemaName, tn.TableName)
					} else {
						fmt.Printf("实体表: %v\n", tn.TableName)
					}
				}
			}

			return false
		},
	}
	stmts, err := parser.Parse(sql)
	if err != nil {
		log.Println(err)
		return
	}

	ok, err := w.Walk(stmts, nil)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(ok)

	fmt.Printf("SQL: %v\n", stmts[0].AST.String())

}
