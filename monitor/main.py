from playwright.sync_api import sync_playwright
from chat import run as chat_run

def main():
    # TODO: 实现请求服务端获取任务逻辑
    # TODO: 实现 Playwright 监控 DeepSeek 逻辑
    # TODO: 实现结果回传逻辑
    print("GEO Monitor started...")
    
    with sync_playwright() as playwright:
        chat_run(playwright)

if __name__ == "__main__":
    main()

