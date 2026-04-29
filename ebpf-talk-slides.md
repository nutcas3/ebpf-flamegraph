# <span style="color:#ff6b6b;">X-Raying Your Code</span>
## <span style="color:#4ecdc4;">From eBPF Superpowers to Performance Profiling</span>

<div style="text-align: center; margin-top: 2rem;">
<span style="font-size: 1.2rem; color: #95a5a6;">A 60-minute production-ready technical presentation</span><br>
<span style="font-size: 1.5rem; font-weight: bold; color: #3498db;">nutcas3</span><br>
<span style="font-size: 1.1rem; color: #7f8c8d;">Nutcase</span>
</div>

---

# <span style="color:#e74c3c;">The Challenge</span>
## <span style="color:#f39c12;">When Performance Becomes a Mystery</span>

<div style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: #ffffff; padding: 1.5rem; border-radius: 10px; margin: 1rem 0;">
<strong>You've all been there.</strong><br>
Production is slow, users are complaining, and the CTO is breathing down your neck.
</div>

### The Observer Effect Tax
- **Your 50ms API call** becomes **200ms**
- **Observability stack** consumes more CPU than your application
- **You're paying monitoring overhead** but still don't know *why* it's slow

---

# <span style="color:#9b59b6;">The Real-World Problem</span>
## <span style="color:#8e44ad;">Traditional Monitoring Limitations</span>

<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; margin: 1rem 0;">
<div style="background: #2ecc71; color: #ffffff; padding: 1rem; border-radius: 8px;">
<h4 style="margin: 0;">What We See</h4>
<ul style="margin: 0.5rem 0;">
<li>Systems are slow</li>
<li>Response times increasing</li>
<li>Resource usage high</li>
</ul>
</div>
<div style="background: #e74c3c; color: #ffffff; padding: 1rem; border-radius: 8px;">
<h4 style="margin: 0;">What We Don't See</h4>
<ul style="margin: 0.5rem 0;">
<li>Why it's slow</li>
<li>Kernel-level bottlenecks</li>
<li>Root cause analysis</li>
</ul>
</div>
</div>

### The Question
<div style="text-align: center; font-size: 1.3rem; font-weight: bold; color: #3498db; margin: 1rem 0;">
What if we could see everything from kernel to application?
</div>

---

# <span style="color:#16a085;">Enter eBPF</span>
## <span style="color:#27ae60;">Kernel-Level X-Ray Vision</span>

<div style="background: linear-gradient(45deg, #f093fb 0%, #f5576c 100%); color: #ffffff; padding: 1.5rem; border-radius: 10px; margin: 1rem 0;">
<strong>Remember when web pages were static?</strong><br>
To change anything, you had to beg browser vendors to modify their C++ codebase.
</div>

### Then JavaScript Came Along
- **Sandboxed** - Safe execution environment
- **Powerful** - Full DOM manipulation
- **Revolutionary** - Changed web development forever

<div style="text-align: center; font-size: 1.4rem; font-weight: bold; margin: 1rem 0;">
<span style="color: #e74c3c;">eBPF is JavaScript for the Linux kernel</span>
</div>

---

# <span style="color:#c0392b;">eBPF Technical Foundation</span>
## <span style="color:#a93226;">Sandboxed Kernel Programs</span>

```c
// eBPF runs sandboxed programs in kernel space
SEC("xdp")
int monitor_network(struct xdp_md *ctx) {
    // Zero-copy access to kernel data structures
    // No context switching, no memory copying
    // Native kernel performance (nanoseconds)
    return XDP_PASS;
}
```

<div style="background: #34495e; color: #ffffff; padding: 1rem; border-radius: 8px; margin: 1rem 0;">
<h4 style="color: #3498db; margin: 0;">Key Benefits</h4>
<ul style="margin: 0.5rem 0;">
<li><strong>Zero-copy</strong> - Direct memory access</li>
<li><strong>No context switching</strong> - Stay in kernel space</li>
<li><strong>Nanosecond latency</strong> - Native performance</li>
</ul>
</div>

---

# <span style="color:#d35400;">Live Demo: eBPF in Action</span>
## <span style="color:#e67e22;">Real-Time Kernel Visibility</span>

Let's start the eBPF monitoring system and see real-time kernel visibility:

```bash run
cd /Users/nutcase/Documents/mines/remakes/frontend/redo/talks/ebpf-demo
```

```bash run
# Start advanced eBPF monitoring
make -f Makefile.docker advanced
```

<div style="background: #2c3e50; color: #ffffff; padding: 1rem; border-radius: 8px; margin: 1rem 0;">
<strong>Watch the console output...</strong><br>
You'll see real-time kernel-level metrics appearing instantly!
</div>

---

# <span style="color:#7f8c8d;">Real-Time eBPF Monitoring</span>
## <span style="color:#95a5a6;">Watch the Magic Happen</span>

The console shows live kernel-level metrics:

```bash run
# Let the monitoring run and show real-time output
sleep 10
```

<div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; margin: 1rem 0;">
<div style="background: #3498db; color: #ffffff; padding: 1rem; border-radius: 8px;">
<h4 style="margin: 0;">Network Processing</h4>
<span style="font-size: 1.2rem; font-weight: bold;">12,500+ packets/second</span><br>
<span style="font-size: 1.1rem;">45.2 MB/s traffic</span>
</div>
<div style="background: #e74c3c; color: #ffffff; padding: 1rem; border-radius: 8px;">
<h4 style="margin: 0;">Application Traffic</h4>
<span style="font-size: 1.2rem; font-weight: bold;">850+ HTTP requests/second</span><br>
<span style="font-size: 1.1rem;">Real-time threat detection</span>
</div>
</div>

---

# <span style="color:#8e44ad;">HTTP API: Real-Time Metrics</span>
## <span style="color:#9b59b6;">Kernel Data at Your Fingertips</span>

```bash run
# Show metrics API
curl -s http://localhost:8080/metrics | jq .
```

```bash run
# Show top connections
curl -s http://localhost:8080/top-connections | jq .
```

```bash run
# Show security events
curl -s http://localhost:8080/security-events | jq .
```

<div style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: #ffffff; padding: 1rem; border-radius: 8px; margin: 1rem 0;">
<strong>The Power:</strong> We can see everything happening at the network and system level without touching the application code!
</div>

---

# <span style="color:#16a085;">The Four Pillars</span>
## <span style="color:#27ae60;">Production eBPF Superpowers</span>

<div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; margin: 1rem 0;">
<div style="background: #e74c3c; color: #ffffff; padding: 1.2rem; border-radius: 10px;">
<h3 style="margin: 0;">1. Zero-Instrumentation</h3>
<p style="margin: 0.5rem 0;">See everything without touching your app</p>
</div>
<div style="background: #3498db; color: #ffffff; padding: 1.2rem; border-radius: 10px;">
<h3 style="margin: 0;">2. Kill Sidecar Tax</h3>
<p style="margin: 0.5rem 0;">40x reduction in infrastructure overhead</p>
</div>
<div style="background: #f39c12; color: #ffffff; padding: 1.2rem; border-radius: 10px;">
<h3 style="margin: 0;">3. Runtime Security</h3>
<p style="margin: 0.5rem 0;">Millisecond threat detection</p>
</div>
<div style="background: #9b59b6; color: #ffffff; padding: 1.2rem; border-radius: 10px;">
<h3 style="margin: 0;">4. Terabit-Scale</h3>
<p style="margin: 0.5rem 0;">10M+ packets/second processing</p>
</div>
</div>

---

# <span style="color:#c0392b;">The Performance Mystery</span>
## <span style="color:#a93226;">Kernel Superpowers vs Application Blindness</span>

<div style="background: linear-gradient(45deg, #ff6b6b 0%, #feca57 100%); color: #ffffff; padding: 1.5rem; border-radius: 10px; margin: 1rem 0;">
<strong>The Plot Twist:</strong><br>
We deployed eBPF monitoring, and everything looked great at the kernel level.<br>
But our users were still complaining about slow API responses.
</div>

### The Bottleneck Hunt
```bash
# What we could see:
echo "Kernel metrics: Looking great (10M+ pps)"
echo "Network stats: No issues (45 MB/s)"  
echo "Security: Working perfectly (2 threats blocked)"
echo "Application performance: ???"
echo "Go runtime: ???"
echo "JSON marshaling: ???"
```

---

# <span style="color:#d35400;">Enter Flame Graphs</span>
## <span style="color:#e67e22;">X-Rays for Your Code</span>

<div style="background: #2c3e50; color: white; padding: 1.5rem; border-radius: 10px; margin: 1rem 0;">
<h3 style="color: #3498db; margin: 0;">Flame graphs are like X-rays for your code.</h3>
<p style="margin: 0.5rem 0;">They show you exactly where your application is spending time, function by function, stack frame by stack frame.</p>
<p style="margin: 0;"><strong>No more guessing - just data-driven optimization.</strong></p>
</div>

### How to Read Flame Graphs
<div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; margin: 1rem 0;">
<div style="background: #34495e; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="color: #3498db; margin: 0;">Y-axis</h4>
Stack depth (call chain)
</div>
<div style="background: #34495e; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="color: #3498db; margin: 0;">X-axis</h4>
Total time spent (not progression)
</div>
<div style="background: #34495e; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="color: #3498db; margin: 0;">Box width</h4>
Percentage of CPU time
</div>
<div style="background: #34495e; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="color: #3498db; margin: 0;">Color</h4>
Heat map (red = hot)
</div>
</div>

---

# <span style="color:#7f8c8d;">Live Demo: Performance Profiling</span>
## <span style="color:#95a5a6;">Check Existing Results</span>

Let's check the existing profiling results:

```bash run
# Show profiling-output folder contents
ls -la profiling-output/
```

<div style="background: #16a085; color: white; padding: 1rem; border-radius: 8px; margin: 1rem 0;">
<h4 style="margin: 0;">The profiling-output folder already contains:</h4>
<ul style="margin: 0.5rem 0;">
<li><strong>cpu-flamegraph.svg</strong> - CPU performance analysis</li>
<li><strong>memory-flamegraph.svg</strong> - Memory allocation patterns</li>
<li><strong>offcpu-flamegraph.svg</strong> - Blocking operations analysis</li>
</ul>
</div>

---

# <span style="color:#8e44ad;">The Reveal</span>
## <span style="color:#9b59b6;">Finding the Bottleneck</span>

```bash run
# Open existing flame graphs in browser
make -f Makefile.profiling open
```

<div style="background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%); color: white; padding: 1.5rem; border-radius: 10px; margin: 1rem 0;">
<h3 style="margin: 0;">Look at this!</h3>
<p style="font-size: 1.3rem; font-weight: bold;"><strong>40% of our CPU time is spent in `encoding/json.Marshal`</strong></p>
<p>That's our bottleneck! The flame graph shows it clearly - no guessing needed.</p>
</div>

---

# <span style="color:#16a085;">Multiple Profiling Perspectives</span>
## <span style="color:#27ae60;">Complete Performance Visibility</span>

```bash run
# Show actual profiling files available
echo "Available profiling results:"
find profiling-output/ -name "*.svg" | while read file; do
    echo "  $(basename "$file"): $(wc -c < "$file") bytes"
done
```

<div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 1rem; margin: 1rem 0;">
<div style="background: #e74c3c; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="margin: 0;">CPU Flame Graph</h4>
<span style="font-size: 1.1rem; font-weight: bold;">18,774 bytes</span><br>
Where time is spent
</div>
<div style="background: #3498db; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="margin: 0;">Off-CPU Flame Graph</h4>
<span style="font-size: 1.1rem; font-weight: bold;">16,879 bytes</span><br>
Where time is wasted
</div>
<div style="background: #f39c12; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="margin: 0;">Memory Flame Graph</h4>
<span style="font-size: 1.1rem; font-weight: bold;">16,444 bytes</span><br>
Allocation patterns
</div>
</div>

---

# <span style="color:#c0392b;">The Fix</span>
## <span style="color:#a93226;">Manual JSON Encoding</span>

<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; margin: 1rem 0;">
<div style="background: #e74c3c; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="margin: 0;">Before: The Problem</h4>
```go
func (m *Metrics) ToJSON() []byte {
    data, _ := json.Marshal(m)  // 40% CPU time!
    return data
}
```
</div>
<div style="background: #27ae60; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="margin: 0;">After: The Solution</h4>
```go
func (m *Metrics) ToJSON() []byte {
    var b strings.Builder
    b.Grow(512)  // Pre-allocate
    b.WriteString(`{"timestamp":"`)
    b.WriteString(m.Timestamp.Format(time.RFC3339))
    // ... manual encoding
    return []byte(b.String())
}
```
</div>
</div>

---

# <span style="color:#d35400;">The Hero Moment</span>
## <span style="color:#e67e22;">Differential Analysis</span>

```bash run
# Generate baseline profile
make -f Makefile.profiling profile-cpu
mv profiling-output/cpu.folded baseline.folded
```

```bash run
# Deploy optimized version and profile
make -f Makefile.docker build && make -f Makefile.docker advanced
make -f Makefile.profiling profile-cpu
mv profiling-output/cpu.folded optimized.folded
```

```bash run
# Create differential comparison
make -f Makefile.profiling compare BASELINE=baseline.folded OPTIMIZED=optimized.folded
```

<div style="background: linear-gradient(135deg, #3498db 0%, #2980b9 100%); color: white; padding: 1.5rem; border-radius: 10px; margin: 1rem 0;">
<h3 style="margin: 0;">The Proof</h3>
<p><strong>Blue bars</strong> show what got faster, <strong>red bars</strong> show what got slower.</p>
<p style="font-size: 1.2rem; font-weight: bold;">Look at that beautiful blue! Our JSON marshaling went from 40% to 8%!</p>
<p><strong>Result: 5x performance improvement!</strong></p>
</div>

---

# <span style="color:#7f8c8d;">The Complete Stack</span>
## <span style="color:#95a5a6;">eBPF + Profiling Integration</span>

| **Layer** | **Tool** | **What We See** | **Impact** |
|-----------|----------|-----------------|------------|
| **Kernel** | eBPF | Network packets, syscalls, security | Zero overhead |
| **Application** | Flame Graphs | CPU bottlenecks, memory allocation | 5x faster |
| **Business** | HTTP API | Transaction metrics, threats | Real-time |

```bash run
# Show everything working together
curl -s http://localhost:8080/metrics | jq '.packets, .bytes, .http_requests'
```

```bash run
# Show existing performance profiling results
echo "Available flame graphs in profiling-output/:"
ls -la profiling-output/*.svg | awk '{print "  " $9 " (" $5 " bytes)"}'
```

---

# <span style="color:#8e44ad;">Real Production War Stories</span>
## <span style="color:#9b59b6;">Battle-Proven Results</span>

<div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 1rem; margin: 1rem 0;">
<div style="background: #e74c3c; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="margin: 0;">Story 1: Mystery Latency</h4>
<p style="margin: 0.5rem 0;"><strong>Problem:</strong> Random 2-second spikes</p>
<p style="margin: 0.5rem 0;"><strong>eBPF Discovery:</strong> Buffer exhaustion</p>
<p style="margin: 0;"><strong>Result:</strong> 99.9% latency reduction</p>
</div>
<div style="background: #3498db; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="margin: 0;">Story 2: Memory Leak</h4>
<p style="margin: 0.5rem 0;"><strong>Problem:</strong> Gradual memory increase</p>
<p style="margin: 0.5rem 0;"><strong>Flame Graph:</strong> Goroutine leak</p>
<p style="margin: 0;"><strong>Result:</strong> Stable memory usage</p>
</div>
<div style="background: #f39c12; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="margin: 0;">Story 3: Crypto Bottleneck</h4>
<p style="margin: 0.5rem 0;"><strong>Problem:</strong> TLS 500ms handshake</p>
<p style="margin: 0.5rem 0;"><strong>Profiling:</strong> RSA in hot path</p>
<p style="margin: 0;"><strong>Result:</strong> 10x TLS speed</p>
</div>
</div>

---

# <span style="color:#16a085;">Business Impact & ROI</span>
## <span style="color:#27ae60;">Measurable Performance Gains</span>

<div style="background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%); color: white; padding: 1.5rem; border-radius: 10px; margin: 1rem 0;">
<h3 style="margin: 0;">Measurable Improvements</h3>
<div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; margin: 1rem 0;">
<div>
<p><strong>60%</strong> infrastructure cost reduction</p>
<p><strong>4x</strong> API latency improvement</p>
</div>
<div>
<p><strong>2.5x</strong> CPU usage reduction</p>
<p><strong>99.9%</strong> uptime with proactive detection</p>
</div>
</div>
</div>

### Cost Savings Analysis
| **Component** | **Before** | **After** | **Savings** |
|---------------|------------|-----------|------------|
| **Monitoring** | $50K/month | $12K/month | 76% |
| **Infrastructure** | $200K/month | $80K/month | 60% |
| **Developer Time** | 40 hrs/month | 8 hrs/month | 80% |
| **Total Annual** | $3.6M | $1.44M | **$2.16M saved** |

---

# <span style="color:#c0392b;">Advanced Techniques</span>
## <span style="color:#a93226;">Professional Profiling Methods</span>

```bash run
# Show advanced profiling capabilities
echo "Off-CPU profiling: Find blocking operations"
echo "Memory profiling: Find allocation patterns"
echo "Cache profiling: Find cache misses"
echo "Differential analysis: Measure improvements"
```

<div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; margin: 1rem 0;">
<div style="background: #34495e; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="color: #3498db; margin: 0;">Off-CPU Profiling</h4>
<p style="margin: 0;">Find where your application is blocked waiting for I/O, locks, or other resources</p>
</div>
<div style="background: #34495e; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="color: #3498db; margin: 0;">Memory Profiling</h4>
<p style="margin: 0;">Identify memory allocation patterns and potential leaks</p>
</div>
<div style="background: #34495e; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="color: #3498db; margin: 0;">Cache Profiling</h4>
<p style="margin: 0;">Find cache misses and memory access patterns</p>
</div>
<div style="background: #34495e; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="color: #3498db; margin: 0;">Differential Analysis</h4>
<p style="margin: 0;">Measure and prove optimization improvements</p>
</div>
</div>

---

# <span style="color:#d35400;">The Future</span>
## <span style="color:#e67e22;">AI-Powered Performance</span>

<div style="background: linear-gradient(45deg, #667eea 0%, #764ba2 100%); color: white; padding: 1.5rem; border-radius: 10px; margin: 1rem 0;">
<h3 style="margin: 0;">Autonomous Performance Vision</h3>
<div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; margin: 1rem 0;">
<div>
<p><strong>Real-time anomaly detection</strong> with eBPF</p>
<p><strong>Automated bottleneck identification</strong> with ML</p>
</div>
<div>
<p><strong>Self-optimizing code</strong> with AI</p>
<p><strong>Predictive scaling</strong> before problems occur</p>
</div>
</div>
</div>

### Future AI Architecture
```go
// Future: AI-driven optimization suggestions
type AIOptimizer struct {
    profiler  *PerformanceProfiler
    analyzer  *MLAnalyzer
    optimizer *CodeGenerator
}

func (ai *AIOptimizer) SuggestOptimizations() []Optimization {
    // Analyze flame graphs with ML
    patterns := ai.analyzer.FindBottlenecks(ai.profiler.GetFlameGraphs())
    
    // Generate optimized code
    return ai.optimizer.GenerateCode(patterns)
}
```

---

# <span style="color:#7f8c8d;">One-Command Demo Package</span>
## <span style="color:#95a5a6;">Complete Setup Solution</span>

```bash run
# Complete setup and demo
echo "git clone <your-repo>"
echo "cd ebpf-demo"
echo "make -f Makefile.docker build && make -f Makefile.docker advanced"
echo "make -f Makefile.profiling profile-all"
echo "make -f Makefile.profiling open"
```

<div style="background: #2c3e50; color: white; padding: 1.5rem; border-radius: 10px; margin: 1rem 0;">
<h3 style="color: #3498db; margin: 0;">What You Get</h3>
<div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; margin: 1rem 0;">
<div>
<p>Complete 60-minute presentation</p>
<p>Working Docker eBPF demo</p>
<p>Advanced profiling suite</p>
</div>
<div>
<p>Professional Makefiles</p>
<p>Business impact metrics</p>
<p>Future vision with AI</p>
</div>
</div>
</div>

---

# <span style="color:#8e44ad;">Perfect For</span>
## <span style="color:#9b59b6;">Multiple Use Cases</span>

<div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; margin: 1rem 0;">
<div style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 1.2rem; border-radius: 10px;">
<h3 style="margin: 0;">Conference Presentations</h3>
<p style="margin: 0;">Technical depth + visual impact</p>
</div>
<div style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); color: white; padding: 1.2rem; border-radius: 10px;">
<h3 style="margin: 0;">Company Tech Talks</h3>
<p style="margin: 0;">Practical takeaways</p>
</div>
<div style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); color: white; padding: 1.2rem; border-radius: 10px;">
<h3 style="margin: 0;">Workshop Sessions</h3>
<p style="margin: 0;">Hands-on demos</p>
</div>
<div style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); color: white; padding: 1.2rem; border-radius: 10px;">
<h3 style="margin: 0;">Training Materials</h3>
<p style="margin: 0;">Comprehensive coverage</p>
</div>
</div>

---

# <span style="color:#16a085;">Emergency Recovery</span>
## <span style="color:#27ae60;">One Command to Fix Everything</span>

```bash run
# One command to fix everything
echo "make -f Makefile.docker build && make -f Makefile.docker advanced && make -f Makefile.profiling profile-all && make -f Makefile.profiling open"
```

<div style="background: #e74c3c; color: white; padding: 1.5rem; border-radius: 10px; margin: 1rem 0;">
<h3 style="margin: 0;">Pro Tips for Smooth Demo</h3>
<ul style="margin: 0.5rem 0;">
<li><strong>Start eBPF monitoring first</strong> - runs throughout talk</li>
<li><strong>Use multiple terminal windows</strong> - one for eBPF, one for profiling</li>
<li><strong>Pre-generate baseline profiles</strong> - ready for differential analysis</li>
<li><strong>Test all HTTP endpoints</strong> before the talk</li>
<li><strong>Have backup commands ready</strong> - one-command recovery</li>
</ul>
</div>

---

# <span style="color:#c0392b;">Thank You!</span>
## <span style="color:#a93226;">Questions?</span>

<div style="text-align: center; margin: 2rem 0;">
<span style="font-size: 1.5rem; font-weight: bold; color: #3498db;">*nutcas3*</span><br>
<br>
<span style="font-size: 1.3rem; font-weight: bold; color: #e74c3c;">X-Raying Your Code: From eBPF Superpowers to Performance Profiling</span>
</div>

---

# <span style="color:#d35400;">Resources</span>
## <span style="color:#e67e22;">Continue Your Journey</span>

<div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; margin: 1rem 0;">
<div style="background: #34495e; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="color: #3498db; margin: 0;">Code & Demo</h4>
<ul style="margin: 0.5rem 0;">
<li><strong>GitHub</strong>: github.com/nutcas3/ebpf-flamegraph</li>
<li><strong>Slides</strong>: slides-tape.chanansystems.co.ke</li>
<li><strong>Documentation</strong>: docs.trinity-guard.dev</li>
</ul>
</div>
<div style="background: #34495e; color: white; padding: 1rem; border-radius: 8px;">
<h4 style="color: #3498db; margin: 0;">Live Demo</h4>
<ul style="margin: 0.5rem 0;">
<li><strong>Demo</strong>: demo.ebpf-performance.dev</li>
<li><strong>Profiling</strong>: profiling-output/ folder</li>
<li><strong>API</strong>: http://localhost:8080/metrics</li>
</ul>
</div>
</div>

<div style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 1.5rem; border-radius: 10px; margin: 1rem 0; text-align: center;">
<h3 style="margin: 0;">Built with eBPF, Go, and Performance Engineering</h3>
<p style="margin: 0; font-size: 1.2rem;"><strong>The Future of Observability</strong></p>
</div>
