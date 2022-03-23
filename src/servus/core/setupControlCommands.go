package core

import "fmt"

func (i *Instance) setupControlCommands() {
	//////////// banhammer.
	// ban.
	i.Control.AddCommand("/ban", func(args []string) {
		if len(args) == 0 {
			i.Control.SendMessage("Usage: /ban ip_address")
			return
		}
		if err := i.Banhammer.Ban(args[0]); err != nil {
			var msg = "[#ERROR] " + err.Error()
			i.Control.SendMessage(msg)
		}
	})

	// unban.
	i.Control.AddCommand("/unban", func(args []string) {
		if len(args) == 0 {
			i.Control.SendMessage("Usage: /unban ip_address")
			return
		}
		if err := i.Banhammer.Unban(args[0]); err != nil {
			var msg = "[#ERROR] " + err.Error()
			i.Control.SendMessage(msg)
		}
	})

	// hooks.
	i.Banhammer.OnBanned(func(ip string) {
		var msg = fmt.Sprintf(`[#BAN] %v`, ip)
		i.Control.SendMessage(msg)
	})
	i.Banhammer.OnUnbanned(func(ip string) {
		var msg = fmt.Sprintf(`[#UNBAN] %v`, ip)
		i.Control.SendMessage(msg)
	})
}
