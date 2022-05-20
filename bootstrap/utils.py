import os, errno, sys, subprocess, shutil, ctypes

def get_absolute_by_relative(path: str) -> str:
    return os.path.abspath(path)

OS_WINDOWS = "NT"
OS_POSIX = "POSIX"
DATA_PATH = get_absolute_by_relative("./data")

def throw_fatal(err: Exception):
    log("bootstrap", err)
    sys.exit(1)

def move_file(fromPath: str, toPath: str):
    try:
        shutil.move(fromPath, toPath)
    except Exception as e:
        throw_fatal(e)

def copy_file(fromPath: str, toPath: str):
    try:
        shutil.copy(fromPath, toPath, follow_symlinks=True)
    except Exception as e:
        throw_fatal(e)

def remove_dir(path: str):
    try:
        shutil.rmtree(path)
    except Exception as e:
        throw_fatal(e)

def remove_file(path: str):
    try:
        os.remove(path)
    except OSError as e:
        if e.errno != errno.ENOENT:
            throw_fatal(e)

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
        throw_fatal(e)

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

def log(who: str, what: str):
    print("[{who}] {what}".format(who=who, what=what))

def get_hosts_path() -> str:
    hosts_path = "/etc/hosts"
    if getOS() == OS_WINDOWS:
        hosts_path = to_unix_path(os.environ["WINDIR"]) + "/System32/drivers/etc/hosts"
    return hosts_path

def is_command_exists(cmd: str):
    from shutil import which
    return which(cmd) is not None

class CommandsStack:
    def __init__(self):
        self.cmdList: list[str] = []

    def addCommand(self, command: str):
        self.cmdList.append(command)

    def runCommands(self):
        for command in self.cmdList:
            run_command(command)