import readline


class StdioTransport:
    def __init__(self):
        readline.parse_and_bind('set editing-mode emacs')

    def read(self):
        return input(">>> ")

    def send(self, text):
        print(text)
