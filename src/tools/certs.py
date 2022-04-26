import utils
import os

def setup():
    utils.run_mkcert_command("-install")

def get():
    CERT_PEM = "cert.pem"
    KEY_PEM = "key.pem"
    # make certs to current dir
    utils.run_mkcert_command(
        '-key-file {KEY_PEM} -cert-file {CERT_PEM} oklookat.ru "*.oklookat.ru" localhost 127.0.0.1 ::1'
        .format(KEY_PEM=KEY_PEM, CERT_PEM=CERT_PEM))

    # copy root cert to devCerts dir
    mkcertFolder = utils.run_mkcert_command('-CAROOT')
    mkcertFolder = mkcertFolder.replace(os.sep, '/')
    ROOTCERT = "{mkcertFolder}/rootCA.pem".format(mkcertFolder=mkcertFolder)
    MOVEROOT = utils.get_absolute_by_relative("./../devCerts/mkcert_root.pem")
    utils.copyFile(ROOTCERT, MOVEROOT)

    # move certs to devCerts dir
    MOVECERT = utils.get_absolute_by_relative("./../devCerts/cert.pem")
    utils.moveFile(CERT_PEM, MOVECERT)
    MOVEKEY = utils.get_absolute_by_relative("./../devCerts/key.pem")
    utils.moveFile(KEY_PEM, MOVEKEY)
    print("done.")

def remove():
    utils.run_mkcert_command("-uninstall")
    mkcertOutput = utils.run_mkcert_command('-CAROOT')
    utils.removeDir(mkcertOutput)
    print("done.")