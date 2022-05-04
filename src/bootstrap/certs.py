import utils, data

def setup():
    utils.run_mkcert_command("-install")

def get():
    # make certs to current dir
    utils.run_mkcert_command(
        '-key-file key.pem -cert-file cert.pem oklookat.ru "*.oklookat.ru" localhost 127.0.0.1 ::1')

    # copy root cert to devCerts dir
    mkcert_root_dir = utils.run_mkcert_command('-CAROOT')
    mkcert_root_dir = utils.to_unix_path(mkcert_root_dir)
    mkcert_root_dir = "{mkcert_root_dir}/rootCA.pem".format(mkcert_root_dir=mkcert_root_dir)

    copy_root_to = utils.CERTS_DEV_DIR + "/mkcert_root.pem"
    utils.copyFile(mkcert_root_dir, copy_root_to)

    # move certs to devCerts dir
    move_cert_to = utils.CERTS_DEV_DIR + "/cert.pem"
    utils.moveFile("cert.pem", move_cert_to)

    move_key_to = utils.CERTS_DEV_DIR + "/key.pem"
    utils.moveFile("key.pem", move_key_to)
    utils.logger("certs", "done")

def remove():
    utils.run_mkcert_command("-uninstall")
    mkcertOutput = utils.run_mkcert_command('-CAROOT')
    utils.removeDir(mkcertOutput)
    utils.logger("certs", "mkcert dir removed")
