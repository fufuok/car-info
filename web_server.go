package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

const html = `<!DOCTYPE html>
<html lang="zh-CN">

<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
	<meta name="renderer" content="webkit">
	<meta content="IE=edge, chrome=1" http-equiv="X-UA-Compatible">
	<title>提取车牌和手机号</title>
	<style>
        body, input, textarea {
            font: 1.2em/1.8 Arial
        }

        table {
            border-collapse: collapse;
            font-family: Arial, sans-serif;
            width: 1200px;
            margin: 0 auto
        }

        caption {
            font-size: larger;
            margin: 1em auto
        }

        th, td {
            padding: .65em
        }

        th, td {
            border-bottom: 1px solid #dedede;
            border-top: 1px solid #dedede;
            text-align: center
        }

        tbody tr:hover {
            background: linear-gradient(#fff, #efefef)
        }

        input {
            cursor: pointer;
        }

        textarea {
            font-size: 1em;
            margin-top: 10px;
            padding: 10px 20px;
            width: 550px;
            height: 600px;
        }
	</style>
	<script type="text/javascript">
        function ajax(opt) {
            opt = opt || {};
            opt.method = (opt.method || 'POST').toUpperCase();
            opt.url = opt.url || '';
            opt.async = opt.async || true;
            opt.data = opt.data || null;
            opt.success = opt.success || function () {
            };
            var xmlHttp = null;
            if (XMLHttpRequest) {
                xmlHttp = new XMLHttpRequest()
            } else {
                xmlHttp = new ActiveXObject('Microsoft.XMLHTTP')
            }
            var params = [];
            for (var key in opt.data) {
                params.push(key + '=' + opt.data[key])
            }
            var postData = params.join('&');
            if (opt.method.toUpperCase() === 'POST') {
                xmlHttp.open(opt.method, opt.url, opt.async);
                xmlHttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded;charset=utf-8');
                xmlHttp.send(postData)
            } else if (opt.method.toUpperCase() === 'GET') {
                xmlHttp.open(opt.method, opt.url + '?' + postData, opt.async);
                xmlHttp.send(null)
            }
            xmlHttp.onreadystatechange = function () {
                if (xmlHttp.readyState === 4 && xmlHttp.status === 200) {
                    opt.success(xmlHttp.responseText)
                }
            }
        }

        function run() {
            ajax({
                url: "/parse", data: {msg: msgform.msg.value}, success: function (r) {
                    msgform.result.value = r
                }
            })
        }
	</script>
</head>

<body>
<form name="msgform">
	<table>
		<caption><input type="button" value="提取车牌和手机号" onclick="run()"></caption>
		<thead>
		<tr>
			<th>源信息</th>
			<th>结果</th>
		</thead>
		<tbody>
		<tr>
			<td>
				<label>
					<textarea placeholder="粘贴源信息, 每一行一条消息" name="msg"></textarea>
				</label>
			</td>
			<td>
				<label>
					<textarea placeholder="提取结果复制到 Excel (Ctrl+A, Ctrl+C)" name="result"></textarea>
				</label>
			</td>
		</tr>
		</tbody>
	</table>
</form>
</body>

</html>`

func initWebServer() error {
	log.Println("Web 服务启动\n浏览器访问: http://127.0.0.1:16888\n关闭窗口或 Ctrl+C 可退出")

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return c.SendString(html)
	})

	app.Post("/parse", func(c *fiber.Ctx) error {
		res := scanMsg(c.FormValue("msg"))
		if res == "" {
			res = "xxx.提取消息结果为空"
		}
		return c.SendString(res)
	})

	return app.Listen(":16888")
}
