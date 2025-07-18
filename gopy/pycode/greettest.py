class Greeter:
    def __init__(self, name):
        self.name = name

    def greet(self, times):
        return f"Hellogg, {self.name}! " * times