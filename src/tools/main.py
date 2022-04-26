import sys
import certs, nginx, utils

def main():
    is_admin = utils.isAdmin()
    if is_admin == False:
        print("run as admin / root")
        sys.exit(1)
    print("""
    -- utils --
    What to do?
    1. [dev certs] Setup
    2. [dev certs] Get & copy
    3. [dev certs] Remove
    4. [nginx] Symlink config
    5. [nginx] Symlink mime
    """)
    what = input("Type a digit: ")
    match what:
        case "1":
            certs.setup()
        case "2":
            certs.get()
        case "3":
            certs.remove()
        case "4":
            nginx.symlink_config()
        case "5":
            nginx.symlink_mime()
        case _:
            main()

# run
main()