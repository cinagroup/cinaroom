#!/usr/bin/env python3
"""
WebSocket 客户端测试脚本
用于测试 WebSocket 终端连接和命令执行
"""

import asyncio
import websockets
import json
import time
import sys
from typing import Optional

class WebSocketTerminalTester:
    def __init__(self, base_url: str = "ws://localhost:8080"):
        self.base_url = base_url
        self.token: Optional[str] = None
        self.session_id: Optional[str] = None
        self.websocket: Optional[websockets.WebSocketClientProtocol] = None
        
    async def connect_terminal(self, token: str) -> bool:
        """连接到终端"""
        self.token = token
        url = f"{self.base_url}/api/v1/ws/terminal?token={token}"
        
        try:
            self.websocket = await websockets.connect(url)
            print(f"✓ Connected to terminal: {url}")
            
            # 等待会话创建消息
            response = await asyncio.wait_for(self.websocket.recv(), timeout=5.0)
            msg = json.loads(response)
            
            if msg.get('type') == 'session_created':
                self.session_id = msg.get('session_id')
                print(f"✓ Session created: {self.session_id}")
                return True
            else:
                print(f"✗ Unexpected response: {msg}")
                return False
                
        except Exception as e:
            print(f"✗ Connection failed: {e}")
            return False
    
    async def send_command(self, command: str) -> str:
        """发送命令"""
        if not self.websocket:
            raise Exception("Not connected")
        
        msg = {
            'type': 'input',
            'payload': command + '\n',
            'timestamp': int(time.time())
        }
        
        await self.websocket.send(json.dumps(msg))
        print(f"→ Sent command: {command}")
        
        # 接收输出
        output = []
        try:
            while True:
                response = await asyncio.wait_for(self.websocket.recv(), timeout=2.0)
                msg = json.loads(response)
                
                if msg.get('type') == 'output':
                    output.append(msg.get('payload', ''))
                    # 如果收到提示符，认为命令完成
                    if '$' in msg.get('payload', '') or '#' in msg.get('payload', ''):
                        break
                elif msg.get('type') == 'heartbeat_ack':
                    continue
                else:
                    output.append(str(msg))
        except asyncio.TimeoutError:
            pass
        
        return ''.join(output)
    
    async def send_heartbeat(self) -> bool:
        """发送心跳"""
        if not self.websocket:
            return False
        
        msg = {
            'type': 'heartbeat',
            'timestamp': int(time.time())
        }
        
        await self.websocket.send(json.dumps(msg))
        
        try:
            response = await asyncio.wait_for(self.websocket.recv(), timeout=5.0)
            msg = json.loads(response)
            return msg.get('type') == 'heartbeat_ack'
        except:
            return False
    
    async def test_commands(self):
        """测试命令执行"""
        commands = [
            'echo "Hello from WebSocket"',
            'pwd',
            'whoami',
            'date',
            'ls -la',
        ]
        
        for cmd in commands:
            print(f"\n{'='*50}")
            output = await self.send_command(cmd)
            print(f"Output:\n{output}")
            await asyncio.sleep(0.5)
    
    async def stress_test(self, concurrent: int = 10, requests: int = 100):
        """压力测试"""
        print(f"\nStarting stress test: {concurrent} concurrent, {requests} requests")
        
        async def worker(worker_id: int):
            success = 0
            failed = 0
            
            for i in range(requests // concurrent):
                try:
                    if await self.send_heartbeat():
                        success += 1
                    else:
                        failed += 1
                except:
                    failed += 1
                
                await asyncio.sleep(0.1)
            
            return success, failed
        
        tasks = [worker(i) for i in range(concurrent)]
        results = await asyncio.gather(*tasks)
        
        total_success = sum(r[0] for r in results)
        total_failed = sum(r[1] for r in results)
        
        print(f"\nStress test completed:")
        print(f"  Success: {total_success}")
        print(f"  Failed: {total_failed}")
        print(f"  Success rate: {total_success / (total_success + total_failed) * 100:.2f}%")
    
    async def close(self):
        """关闭连接"""
        if self.websocket:
            await self.websocket.close()
            print("✓ Connection closed")

async def main():
    """主函数"""
    import argparse
    
    parser = argparse.ArgumentParser(description='WebSocket Terminal Tester')
    parser.add_argument('--url', default='ws://localhost:8080', help='WebSocket URL')
    parser.add_argument('--token', default='test-token', help='Auth token')
    parser.add_argument('--stress', action='store_true', help='Run stress test')
    parser.add_argument('--concurrent', type=int, default=10, help='Concurrent connections')
    parser.add_argument('--requests', type=int, default=100, help='Total requests')
    
    args = parser.parse_args()
    
    tester = WebSocketTerminalTester(args.url)
    
    try:
        # 连接
        if await tester.connect_terminal(args.token):
            # 测试命令
            await tester.test_commands()
            
            # 压力测试
            if args.stress:
                await tester.stress_test(args.concurrent, args.requests)
    except KeyboardInterrupt:
        print("\nTest interrupted")
    except Exception as e:
        print(f"Test failed: {e}")
    finally:
        await tester.close()

if __name__ == '__main__':
    asyncio.run(main())
