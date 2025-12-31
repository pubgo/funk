package featurehttp

const htmlTemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Feature Flags 管理</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 12px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        
        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
        }
        
        .header p {
            opacity: 0.9;
            font-size: 1.1em;
        }
        
        .content {
            padding: 30px;
        }
        
        .search-box {
            margin-bottom: 30px;
        }
        
        .search-box input {
            width: 100%;
            padding: 12px 20px;
            font-size: 16px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            transition: border-color 0.3s;
        }
        
        .search-box input:focus {
            outline: none;
            border-color: #667eea;
        }
        
        .feature-list {
            display: grid;
            gap: 20px;
        }
        
        .feature-card {
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            padding: 20px;
            transition: all 0.3s;
            background: #fafafa;
        }
        
        .feature-card:hover {
            border-color: #667eea;
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
            transform: translateY(-2px);
        }
        
        .feature-header {
            display: flex;
            justify-content: space-between;
            align-items: start;
            margin-bottom: 15px;
        }
        
        .feature-name {
            font-size: 1.3em;
            font-weight: bold;
            color: #333;
        }
        
        .feature-badges {
            display: flex;
            gap: 8px;
            flex-wrap: wrap;
        }
        
        .badge {
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 0.85em;
            font-weight: 500;
        }
        
        .badge-type {
            background: #667eea;
            color: white;
        }
        
        .badge-deprecated {
            background: #f44336;
            color: white;
        }
        
        .badge-sensitive {
            background: #ff9800;
            color: white;
        }
        
        .badge-immutable {
            background: #9e9e9e;
            color: white;
        }
        
        .feature-usage {
            color: #666;
            margin-bottom: 15px;
            font-style: italic;
        }
        
        .feature-value {
            display: flex;
            gap: 10px;
            align-items: center;
            margin-bottom: 15px;
        }
        
        .value-display {
            flex: 1;
            padding: 10px 15px;
            background: white;
            border: 2px solid #e0e0e0;
            border-radius: 6px;
            font-family: 'Courier New', monospace;
            font-size: 14px;
            word-break: break-all;
        }
        
        .value-input {
            flex: 1;
            padding: 10px 15px;
            border: 2px solid #e0e0e0;
            border-radius: 6px;
            font-family: 'Courier New', monospace;
            font-size: 14px;
            display: none;
        }
        
        .value-input:focus {
            outline: none;
            border-color: #667eea;
        }
        
        .btn {
            padding: 10px 20px;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            transition: all 0.3s;
        }
        
        .btn-edit {
            background: #667eea;
            color: white;
        }
        
        .btn-edit:hover {
            background: #5568d3;
        }
        
        .btn-save {
            background: #4caf50;
            color: white;
            display: none;
        }
        
        .btn-save:hover {
            background: #45a049;
        }
        
        .btn-cancel {
            background: #9e9e9e;
            color: white;
            display: none;
        }
        
        .btn-cancel:hover {
            background: #757575;
        }
        
        .btn:disabled {
            opacity: 0.5;
            cursor: not-allowed;
        }
        
        .feature-tags {
            margin-top: 15px;
            padding-top: 15px;
            border-top: 1px solid #e0e0e0;
        }
        
        .tags-label {
            font-size: 0.9em;
            color: #999;
            margin-bottom: 8px;
        }
        
        .tags-content {
            display: flex;
            flex-wrap: wrap;
            gap: 8px;
        }
        
        .tag-item {
            padding: 4px 10px;
            background: #e3f2fd;
            border-radius: 4px;
            font-size: 0.85em;
            color: #1976d2;
        }
        
        .message {
            padding: 15px;
            border-radius: 6px;
            margin-bottom: 20px;
            display: none;
        }
        
        .message-success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        
        .message-error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        
        .empty-state {
            text-align: center;
            padding: 60px 20px;
            color: #999;
        }
        
        .empty-state h3 {
            font-size: 1.5em;
            margin-bottom: 10px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Feature Flags 管理</h1>
            <p>查看和管理应用程序的功能标志</p>
        </div>
        
        <div class="content">
            <div id="message" class="message"></div>
            
            <div class="search-box">
                <input type="text" id="searchInput" placeholder="搜索功能标志..." oninput="filterFeatures()">
            </div>
            
            <div class="feature-list" id="featureList">
                {{range .Flags}}
                <div class="feature-card" data-name="{{.Name}}" data-usage="{{.Usage}}">
                    <div class="feature-header">
                        <div class="feature-name">{{.Name}}</div>
                        <div class="feature-badges">
                            <span class="badge badge-type">{{.Type}}</span>
                            {{if .Deprecated}}
                            <span class="badge badge-deprecated">已废弃</span>
                            {{end}}
                            {{if .Sensitive}}
                            <span class="badge badge-sensitive">敏感</span>
                            {{end}}
                            {{if not .Mutable}}
                            <span class="badge badge-immutable">不可修改</span>
                            {{end}}
                        </div>
                    </div>
                    
                    <div class="feature-usage">{{.Usage}}</div>
                    
                    <div class="feature-value">
                        <div class="value-display" id="display-{{.Name}}">{{.ValueString}}</div>
                        <input type="text" class="value-input" id="input-{{.Name}}" value="{{.ValueString}}" data-original="{{.ValueString}}">
                        <button class="btn btn-edit" onclick="editFeature('{{.Name}}')" id="edit-{{.Name}}" {{if not .Mutable}}disabled{{end}}>编辑</button>
                        <button class="btn btn-save" onclick="saveFeature('{{.Name}}')" id="save-{{.Name}}">保存</button>
                        <button class="btn btn-cancel" onclick="cancelEdit('{{.Name}}')" id="cancel-{{.Name}}">取消</button>
                    </div>
                    
                    {{if .Tags}}
                    <div class="feature-tags">
                        <div class="tags-label">标签:</div>
                        <div class="tags-content">
                            {{range $key, $value := .Tags}}
                            <span class="tag-item">{{$key}}: {{$value}}</span>
                            {{end}}
                        </div>
                    </div>
                    {{end}}
                </div>
                {{end}}
            </div>
            
            {{if not .Flags}}
            <div class="empty-state">
                <h3>暂无功能标志</h3>
                <p>还没有注册任何功能标志</p>
            </div>
            {{end}}
        </div>
    </div>
    
    <script>
        const prefix = '{{.Prefix}}';
        
        function filterFeatures() {
            const searchInput = document.getElementById('searchInput');
            const filter = searchInput.value.toLowerCase();
            const cards = document.querySelectorAll('.feature-card');
            
            cards.forEach(card => {
                const name = card.getAttribute('data-name').toLowerCase();
                const usage = card.getAttribute('data-usage').toLowerCase();
                if (name.includes(filter) || usage.includes(filter)) {
                    card.style.display = '';
                } else {
                    card.style.display = 'none';
                }
            });
        }
        
        function editFeature(name) {
            const display = document.getElementById('display-' + name);
            const input = document.getElementById('input-' + name);
            const editBtn = document.getElementById('edit-' + name);
            const saveBtn = document.getElementById('save-' + name);
            const cancelBtn = document.getElementById('cancel-' + name);
            
            display.style.display = 'none';
            input.style.display = 'block';
            editBtn.style.display = 'none';
            saveBtn.style.display = 'inline-block';
            cancelBtn.style.display = 'inline-block';
            
            input.focus();
            input.select();
        }
        
        function cancelEdit(name) {
            const display = document.getElementById('display-' + name);
            const input = document.getElementById('input-' + name);
            const editBtn = document.getElementById('edit-' + name);
            const saveBtn = document.getElementById('save-' + name);
            const cancelBtn = document.getElementById('cancel-' + name);
            const original = input.getAttribute('data-original');
            
            input.value = original;
            display.style.display = 'block';
            input.style.display = 'none';
            editBtn.style.display = 'inline-block';
            saveBtn.style.display = 'none';
            cancelBtn.style.display = 'none';
        }
        
        async function saveFeature(name) {
            const input = document.getElementById('input-' + name);
            const value = input.value;
            
            try {
                const response = await fetch(prefix + '/api/update', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        name: name,
                        value: value
                    })
                });
                
                const data = await response.json();
                
                if (response.ok && data.success) {
                    const display = document.getElementById('display-' + name);
                    display.textContent = data.flag.valueString;
                    input.value = data.flag.valueString;
                    input.setAttribute('data-original', data.flag.valueString);
                    
                    showMessage('更新成功: ' + data.message, 'success');
                    cancelEdit(name);
                } else {
                    showMessage('更新失败: ' + (data.message || response.statusText), 'error');
                }
            } catch (error) {
                showMessage('更新失败: ' + error.message, 'error');
            }
        }
        
        function showMessage(text, type) {
            const messageDiv = document.getElementById('message');
            messageDiv.textContent = text;
            messageDiv.className = 'message message-' + type;
            messageDiv.style.display = 'block';
            
            setTimeout(() => {
                messageDiv.style.display = 'none';
            }, 3000);
        }
        
        // 支持 Enter 键保存
        document.addEventListener('keydown', function(e) {
            if (e.key === 'Enter') {
                const activeInput = document.activeElement;
                if (activeInput && activeInput.classList.contains('value-input')) {
                    const name = activeInput.id.replace('input-', '');
                    const saveBtn = document.getElementById('save-' + name);
                    if (saveBtn && saveBtn.style.display !== 'none') {
                        saveFeature(name);
                    }
                }
            }
        });
    </script>
</body>
</html>`
