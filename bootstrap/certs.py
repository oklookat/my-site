import utils

CERTS_DEV_DIR = "../nginx/.devcontainer/certs"

def run_mkcert_command(command: str) -> str:
    if not utils.is_command_exists("mkcert"):
        utils.throw_fatal(Exception("mkcert not installed"))
    try:
        out = utils.run_command("mkcert " + command, True)
        return out.strip()
    except Exception as e:
        utils.throw_fatal(e)

def setup():
    run_mkcert_command("-install")

def get():
    # make certs to current dir
    run_mkcert_command(
        '-key-file key.pem -cert-file cert.pem oklookat.ru "*.oklookat.ru" localhost 127.0.0.1 ::1')

    # copy root cert to devCerts dir
    mkcert_root_dir = run_mkcert_command('-CAROOT')
    mkcert_root_dir = utils.to_unix_path(mkcert_root_dir)
    mkcert_root_dir = "{mkcert_root_dir}/rootCA.pem".format(mkcert_root_dir=mkcert_root_dir)

    copy_root_to = CERTS_DEV_DIR + "/mkcert_root.pem"
    utils.copy_file(mkcert_root_dir, copy_root_to)

    # move certs to devCerts dir
    move_cert_to = CERTS_DEV_DIR + "/cert.pem"
    utils.move_file("cert.pem", move_cert_to)

    move_key_to = CERTS_DEV_DIR + "/key.pem"
    utils.move_file("key.pem", move_key_to)
    utils.log("certs", "done")

def remove():
    run_mkcert_command("-uninstall")
    mkcert_output = run_mkcert_command('-CAROOT')
    utils.remove_dir(mkcert_output)
    utils.log("certs", "removed")
