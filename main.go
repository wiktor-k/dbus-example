package main

import (
	"fmt"
	"os"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

const intro = `
<node>
	<interface name="biz.metacode.Demo">
		<method name="Foo">
			<arg direction="out" type="s"/>
		</method>
		<method name="Bar">
			<arg direction="in" type="s"/>
			<arg direction="out" type="s"/>
		</method>
	</interface>` + introspect.IntrospectDataString + `</node> `

type foo string

func (f foo) Foo() (string, *dbus.Error) {
	fmt.Println(f)
	return string(f), nil
}

func (f foo) Bar(x string) (string, *dbus.Error) {
	fmt.Println(f)
	return string(x + string(f) + x), nil
}

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	f := foo("Testing DBus!")
	conn.Export(f, "/biz/metacode/Example", "biz.metacode.Demo")
	conn.Export(introspect.Introspectable(intro), "/biz/metacode/Example",
		"org.freedesktop.DBus.Introspectable")

	reply, err := conn.RequestName("biz.metacode.Demo", dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}
	fmt.Println("Listening on biz.metacode.Demo / /biz/metacode/Example ...")
	select {}
}
