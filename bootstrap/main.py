import sys
import certs, utils, hosts, ancient_ui

def main():
    if not utils.isAdmin():
        utils.throw_fatal(Exception("run as admin/root"))

    d_exit = ancient_ui.ChooserCommand("Exit", lambda: sys.exit(0))
    dev_mkcert_setup = ancient_ui.ChooserCommand("[dev/certs] Setup", certs.setup)
    dev_mkcert_gen_copy = ancient_ui.ChooserCommand("[dev/certs] Generate & Copy to devCerts", certs.get)
    dev_mkcert_remove = ancient_ui.ChooserCommand("[dev/certs] Remove", certs.remove)
    dev_hosts_add = ancient_ui.ChooserCommand("[dev/hosts] Add/remove hosts", hosts.add)

    chooser = ancient_ui.Chooser(d_exit, dev_mkcert_setup, dev_mkcert_gen_copy, 
    dev_mkcert_remove, dev_hosts_add)
    chooser.run()

# run
main()
