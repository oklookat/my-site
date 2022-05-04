import utils, data

SYSTEM_HOSTS = utils.get_hosts_path()

def add():
    # data hosts
    paste_lines = []
    with open(utils.HOSTS_PATH, "r") as hostsToPaste:
        paste_lines = hostsToPaste.readlines()

    # system hosts
    hosts_lines = []
    with open(SYSTEM_HOSTS, "r") as hosts:
        hosts_lines = hosts.readlines()


    final_lines = []
    for paste in paste_lines:
        paste_line = paste.strip() + "\n"
        line_exists = False
        for line in hosts_lines:
            line_strip = line.strip() + "\n"
            line_exists = line_strip in paste_line
            if line_exists:
                break
        if not line_exists or len(hosts_lines) < 1:
            final_lines.append(paste_line)

    if len(final_lines) < 1:
        remove()
        return

    hosts = open(SYSTEM_HOSTS, "a")
    hosts.write("\n\n")
    hosts.writelines(final_lines)
    hosts.close()
    utils.logger("hosts", "added")


def remove():
    # data hosts
    exclude_lines = []
    with open(utils.HOSTS_PATH, "r") as exclude:
        exclude_lines = exclude.readlines()

    # system hosts
    hosts_lines = []
    with open(SYSTEM_HOSTS, "r") as hosts:
        hosts_lines = hosts.readlines()

    final_lines = []
    for line in hosts_lines:
        line_strip = line.strip() + "\n"
        need_to_exclude = False
        for exclude in exclude_lines:
            exclude_strip = exclude.strip() + "\n"
            need_to_exclude = line_strip in exclude_strip
            if need_to_exclude:
                break
        if not need_to_exclude:
            final_lines.append(line)

    hosts = open(SYSTEM_HOSTS, "w")
    hosts.writelines(final_lines)
    hosts.close()
    utils.logger("hosts", "removed")
