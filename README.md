
![ChatGPT Image 13 de abr  de 2025, 00_18_33](https://github.com/user-attachments/assets/90a7c815-7e5a-4cd9-aa74-52412eb5ebd2)


🔍 kube-validator(Under Construction)

`kube-validator` is a Go-based CLI tool designed to validate Kubernetes manifests against best practices and custom policies. It helps platform teams, DevOps engineers, and developers ensure correctness, consistency, and quality before shipping YAMLs into clusters.

---

## 🚀 Features

- **⚙️ Custom Policy Validation**: Write and enforce your own validation rules for Kubernetes objects.
- **🧠 Semantic Awareness**: Validates fields based on expected types, required presence, and Kubernetes API semantics.
- **📦 Multi-Format Support**: Accepts both single YAML files and directories with multi-doc YAMLs.
- **📣 Human-Friendly Output**: Results are color-coded, categorized (`Error`, `Warning`, `Suggestion`), and explained.
- **🧪 Testable & CI/CD Friendly**: Can run in pipelines or pre-commit hooks, with optional JSON output.

---

## 📦 Installation

### Prerequisites
- Go 1.18 or later

### Build
```bash
git clone https://github.com/giovanni-gava/kube-validator.git
cd kube-validator
go build -o kube-validator main.go
```

### Run
```bash
./kube-validator validate -f ./examples/deployment.yaml
```

---

## 🛠️ Usage

```bash
kube-validator validate [flags]
```

### Common Flags
| Flag        | Description                                  |
|-------------|----------------------------------------------|
| `-f, --file` | Path to manifest file or directory            |
| `--strict`   | Treat warnings as errors                      |
| `--output`   | Output format: `pretty`, `json`              |

### Example
```bash
kube-validator validate -f ./examples/ --output pretty
```

---

## 📁 Project Structure

```bash
kube-validator/
├── cmd/                # CLI command parsing
├── internal/           # Rule engines and validators
├── pkg/                # Shared utilities and types
├── examples/           # Sample YAMLs for testing
├── main.go             # Entry point
├── go.mod              # Go module definition
└── README.md
```

---

## 🔬 Validator Engine

Each validator is modular and independent. For example:
- `metadata.name` presence and format
- `spec.replicas` must be a positive integer
- `resources.requests` and `limits` must be defined
- `livenessProbe` and `readinessProbe` should be configured

> You can easily add your own rules under `internal/validators/`.

---

## 🧪 Testing

```bash
go test ./...
```
Unit tests are powered by [Testify](https://github.com/stretchr/testify) and cover:
- Manifest parsing
- Rule engine
- Output rendering

---

## 🧬 Output Example

```yaml
[Error] Missing field: spec.replicas
  → Deployment should declare how many replicas it expects.
  → Docs: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/

[Warning] No resource limits defined
  → It is recommended to set CPU and memory limits.

[Suggestion] Consider using livenessProbe
  → This helps Kubernetes automatically detect unhealthy containers.
```

---

## 🤝 Contributing

1. Fork it 💡
2. Create your feature branch (`git checkout -b feat/new-rule`)
3. Commit your changes (`git commit -am 'add cool rule'`)
4. Push to the branch (`git push origin feat/new-rule`)
5. Create a Pull Request 🚀

---

## 📄 License

MIT License © [Giovanni Gava](https://github.com/giovanni-gava)

---

## 📸 Logo
Logo is generated below ⬇️ — feel free to use it as the project badge!
