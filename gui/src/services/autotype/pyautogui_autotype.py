import pyautogui

def autotype_pyautogui(text):
    print("autotypeing with pyautogui")
    pyautogui.write(text, interval=0.02)