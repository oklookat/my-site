import sys
import certs, nginx, utils, hosts


def main():
    is_admin = utils.isAdmin()
    if is_admin == False:
        utils.logger("bootstrap", "run as admin/root")
        sys.exit(1)

    print(
        """
-------- oklookat-site bootstrap --------
-- development
1. [dev/certs] Setup (mkcert)
2. [dev/certs] Generate & Copy to devCerts
3. [dev/certs] Remove
4. [dev/hosts] Add/remove site hosts
5. [dev/nginx] Symlink config
6. [dev/nginx] Symlink mime types
----------------
    """
    )

    what = input("Type a digit: ")
    try:
        match what:
            case "1":
                certs.setup()
            case "2":
                certs.get()
            case "3":
                certs.remove()
            case "4":
                hosts.add()
            case "5":
                nginx.symlink_dev_config()
            case "6":
                nginx.symlink_dev_mime()
            case _:
                main()
    except Exception as e:
        utils.logger("bootstrap", e)


# run
main()
