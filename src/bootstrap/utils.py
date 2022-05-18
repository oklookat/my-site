import os, errno, sys, subprocess, shutil, ctypes

def get_absolute_by_relative(path: str) -> str:
    return os.path.abspath(path)

OS_WINDOWS = "NT"
OS_POSIX = "POSIX"
DATA_PATH = get_absolute_by_relative("./data")
CONFIG_PATH = get_absolute_by_relative("./config.yaml")
HOSTS_PATH = DATA_PATH + "/hosts.txt"
CERTS_DEV_DIR = DATA_PATH + "/devCerts"
NGINX_MIME_TYPES_PATH = DATA_PATH + "/mime.types"
NGINX_DEV_CONF_PATH = DATA_PATH + "/nginx.dev.conf"


def throwFatalErr(err: Exception):
    print(
        """
        -------- ERROR --------
        {err}
        -------- ERROR --------
        """.format(
            err=err
        )
    )
    sys.exit(1)


def moveFile(fromPath: str, toPath: str):
    try:
        shutil.move(fromPath, toPath)
    except Exception as e:
        throwFatalErr(e)


def copyFile(fromPath: str, toPath: str):
    try:
        shutil.copy(fromPath, toPath, follow_symlinks=True)
    except Exception as e:
        throwFatalErr(e)


def removeDir(path: str):
    try:
        shutil.rmtree(path)
    except Exception as e:
        throwFatalErr(e)


def removeFile(path: str):
    try:
        os.remove(path)
    except OSError as e:
        if e.errno != errno.ENOENT:
            throwFatalErr(e)


def getOS() -> str:
    return os.name.upper()


def run_command(command: str, exit_if_error: bool = True) -> str:
    try:
        out = subprocess.run(command, shell=True, check=True, capture_output=True)
        res = out.stdout.decode("UTF-8")
        return res
    except Exception as e:
        if exit_if_error == False:
            raise e
        throwFatalErr(e)


def run_mkcert_command(command: str) -> str:
    mkcert_path = "mkcert "
    try:
        out = run_command(mkcert_path + command, True)
        return out.strip()
    except Exception as e:
        logger("mkcert", e)
        logger("mkcert", "mkcert not installed (?)")

def to_unix_path(val: str) -> str:
    return val.replace(os.sep, '/')

# get main.py dir
def get_execution_dir() -> str:
    cwd = to_unix_path(os.getcwd())
    if not cwd.endswith("/"):
        cwd += "/"
    cwd += "main.py"
    return cwd


def isAdmin() -> bool:
    is_admin = False
    try:
        if getOS() == OS_WINDOWS:
            is_admin = ctypes.windll.shell32.IsUserAnAdmin() != 0
        else:
            is_admin = os.getuid() == 0
    except Exception as _:
        raise Exception("Unsupported OS.")
    return is_admin

class CommandsStack:
    def __init__(self):
        self.cmdList: list[str] = []

    def addCommand(self, command: str):
        self.cmdList.append(command)

    def runCommands(self):
        for command in self.cmdList:
            run_command(command)


def get_hosts_path() -> str:
    hostsPath = "/etc/hosts"
    if getOS() == OS_WINDOWS:
        hostsPath = to_unix_path(os.environ["WINDIR"]) + "/System32/drivers/etc/hosts"
    return hostsPath

def logger(who: str, what: str):
    print("[{who}] {what}.".format(who=who, what=what))

