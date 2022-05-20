import typing

COLOR_RESET = "\u001b[0m"
COLOR_BLUE = "\u001b[34;1m"


class ChooserCommand:
    def __init__(self, description: str, callback: typing.Callable) -> None:
        self.description: str = description
        self.callback: typing.Callable = callback
    def run_callback(self) -> None:
        self.callback()

class Chooser:
    def __init__(self, *cmds: ChooserCommand) -> None:
        self.is_wrong_digit = False
        self.uid: int = 0
        self.commands: dict[int, ChooserCommand] = dict()
        if len(cmds) > 0:
            for command in cmds:
                self.add(command)           
    def add(self, cmd: ChooserCommand):
        self.uid += 1
        self.commands[self.uid] = cmd

    def run(self):
        result = False
        while not result:
            result = self.__run()
            
    def __run(self) -> bool:
        if not self.is_wrong_digit:
            print("------------")
            for uid in self.commands:
                command = self.commands[uid]
                print("{color}{uid}{reset} | {desc}"
                .format(color=COLOR_BLUE, reset=COLOR_RESET, uid=uid, desc=command.description))
            print("------------")
        else:
            print("Choose correct digit.")

        result = input("Digit: ")
        result_uid = -1
        is_wrong = False

        try:
            result_uid = int(result)
        except Exception as _:
            is_wrong = True
        if result_uid not in self.commands:
            is_wrong = True
        if is_wrong:
            self.is_wrong_digit = True
            return False

        command = self.commands[result_uid]
        command.callback()
        self.is_wrong_digit = False
        return True