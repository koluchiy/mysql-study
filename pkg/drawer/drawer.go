package drawer

import (
	"github.com/koluchiy/mysql-study/pkg/connection"
	"fmt"
	"strings"
)

type Message struct {
	Msg string
}

func (m Message) IsEmpty() bool {
	return len(m.Msg) == 0
}

func NewMessageString(msg string) Message {
	return Message{Msg: msg}
}

func (m Message) String() string {
	return m.Msg
}

type Drawer interface {
	DrawQueryResult(data connection.QueryResult)
	DrawMessage(m Message)
}

type drawer struct {

}

func NewDrawer() *drawer {
	instance := &drawer{}

	return instance
}

func (v *drawer) DrawQueryResult(data connection.QueryResult) {
	fmt.Println("| " + strings.Join(data.Columns, " | ") + " |")

	for _, row := range data.Rows {
		strs := make([]string, len(row))

		for i, c := range row {
			strs[i] = string(c)
		}

		fmt.Println("| " + strings.Join(strs, " | ") + " |")
	}
}

func (v *drawer) DrawMessage(m Message) {
	fmt.Println(m)
}