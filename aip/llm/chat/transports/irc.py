class IrcTransport:
    def __init__(self):
        pass

    def read(self):
        return input(">>> ")

    def send(self, text):
        print(text)
