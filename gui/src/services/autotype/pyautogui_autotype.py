import pyautogui

def autotype_pyautogui(text):
    print("autotypeing with pyautogui")
    pyautogui.write(text, interval=0.02)

if __name__ == "__main__":
    autotype_pyautogui("hello world")