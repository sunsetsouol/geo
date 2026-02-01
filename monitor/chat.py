import re
import time
import json
import os
import requests
from playwright.sync_api import Playwright, sync_playwright, expect

SERVER_URL = "http://localhost:8080"

def fetch_pending_tasks():
    try:
        response = requests.get(f"{SERVER_URL}/api/tasks/pending")
        if response.status_code == 200:
            return response.json()
        else:
            print(f"Failed to fetch tasks: {response.status_code}")
            return []
    except Exception as e:
        print(f"Error fetching tasks: {e}")
        return []

def upload_task_result(task_id, result_data):
    try:
        response = requests.post(f"{SERVER_URL}/api/tasks/{task_id}/result", json=result_data)
        if response.status_code == 200:
            print(f"Successfully uploaded result for task {task_id}")
        else:
            print(f"Failed to upload result for task {task_id}: {response.status_code}")
    except Exception as e:
        print(f"Error uploading result: {e}")

def run_task(page, task):
    prompt_content = task.get("prompt", {}).get("content", "")
    if not prompt_content:
        print(f"No prompt content for task {task['id']}")
        return None

    print(f"Running task {task['id']} with prompt: {prompt_content}")
    
    # 模拟 DeepSeek 交互
    page.goto("https://chat.deepseek.com/")
    
    # 确保联网搜索已开启 (根据现有代码逻辑)
    try:
        page.get_by_role("button", name="联网搜索").click()
    except:
        pass # 如果已经开启或找不到按钮，跳过

    page.get_by_role("textbox", name="给 DeepSeek 发送消息").fill(prompt_content)
    page.get_by_role("button").nth(4).click() # 发送按钮

    # 等待回复生成完成 (这里需要根据实际页面结构优化等待逻辑)
    time.sleep(15) 

    # 提取回复内容 (通过点击复制按钮获取剪贴板内容)
    try:
        # 点击复制按钮
        page.get_by_role("button").nth(4).click()
        page.locator(".ds-flex._965abe9 > div").first.click()
        # 获取剪贴板内容
        response_text = page.evaluate("navigator.clipboard.readText()")
        print(response_text)
    except Exception as e:
        print(f"Failed to copy response text: {e}")
        # 降级方案：直接获取 inner_text
        response_text = page.locator(".ds-markdown").last.inner_text()

    # 提取引用链接 (示例逻辑)
    citations = []
    citation_elements = page.locator("a[href^='http']").all()
    for el in citation_elements:
        href = el.get_attribute("href")
        if href and "deepseek.com" not in href:
            citations.append({"url": href, "title": el.inner_text()})

    return {
        "status": "completed",
        "response_text": response_text,
        "brand_score": 0.0, # 初始评分，后续由分析逻辑填充
        "analysis_report": "{}",
        "citations": citations[:10] # 限制数量
    }

def run(playwright: Playwright) -> None:
    browser = playwright.chromium.launch(headless=False)
    auth_file = os.path.join(os.path.dirname(__file__), "auth.json")
    
    storage_state = None
    if os.path.exists(auth_file):
        with open(auth_file) as f:
            storage_state = json.load(f)
    
    context = browser.new_context(
        locale="zh-CN",
        storage_state=storage_state,
        permissions=["clipboard-read", "clipboard-write"]
    )

    page = context.new_page()

    while True:
        print(f"[{time.strftime('%Y-%m-%d %H:%M:%S')}] Fetching pending tasks...")
        tasks = fetch_pending_tasks()
        
        if not tasks:
            print("No pending tasks found. Sleeping for 60 seconds...")
            time.sleep(60)
            continue

        print(f"Found {len(tasks)} tasks. Starting processing...")
        for task in tasks:
            try:
                result = run_task(page, task)
                if result:
                    upload_task_result(task["id"], result)
                else:
                    upload_task_result(task["id"], {"status": "failed"})
            except Exception as e:
                print(f"Error processing task {task['id']}: {e}")
                upload_task_result(task["id"], {"status": "failed"})
            
            # 任务之间稍微停顿，避免请求过快
            time.sleep(5)

        # 每轮任务完成后保存一次登录状态
        context.storage_state(path=auth_file)
        print("Batch completed. Checking for more tasks...")

    # 理论上不会执行到这里
    context.close()
    browser.close()
