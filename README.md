# Vibe Check

> Simple Git workflow made easy - no more complex commands!

**Vibe Check** is a lightweight CLI tool that simplifies Git version management. Tired of complicated Git commands with endless options? This tool gives you a clean, intuitive interface for creating checkpoints, switching between versions, and cleaning up your commit history.

**Why Vibe Check?**
- Git commands are complex and overwhelming with too many options
- Easy to make mistakes with `git reset`, `git rebase`, `git cherry-pick`
- Perfect for iterative development, experimentation, and AI-assisted coding
- Simplifies the entire Git workflow into simple menu choices

## ✨ Features

### 🎯 **Simple Checkpoint System**
Instead of complex Git commands, just use simple menu options:
- **Create Checkpoint** - Save your current work instantly (like `git add . && git commit`)
- **Change Checkpoint** - Switch between any version (replaces complex `git checkout`/`git reset`)
- **Finalize and Push** - Clean up and push to remote (handles `git rebase`, `git push --force-with-lease`)

### 🚀 **No More Command Hell**
Forget about:
- `git rebase -i HEAD~5` - Complex interactive rebasing
- `git reset --hard HEAD~3` - Dangerous hard resets  
- `git cherry-pick abc123` - Confusing commit picking
- `git push --force-with-lease` - Scary force pushing

## 🔒 Safety & Privacy

### 🛡️ **How It Works**
- **Simple wrapper** - Vibe Check only executes `git` commands in your terminal
- **No network access** - Just runs `git add`, `git commit`, `git checkout`, `git push`, etc.
- **Transparent operations** - Everything it does, you could do manually with Git commands

**Vibe Check handles it all with simple menu navigation!**

## 🛠️ Installation

### NPM (Recommended)
```bash
# Install globally
npm install -g @bomoge/vibe-check

# Use immediately
vibe-check
```

### From Binary
```bash
# Install globally (macOS with Homebrew)
cp vibe-check /opt/homebrew/bin/vibe-check

# Or to /usr/local/bin (requires sudo)
sudo cp vibe-check /usr/local/bin/vibe-check
```

### From Source
```bash
git clone https://github.com/pr0d5h381/vibe-check.git
cd vibe-check
go build -o vibe-check .
```

## 🚦 Usage

### Interactive Mode (TUI)
Navigate to any Git repository and run:

```bash
vibe-check
```

Use arrow keys to navigate the menu, Enter to select, and follow the intuitive interface.

### Command Line Mode
For quick operations and scripting:

```bash
# Create checkpoint with custom note
vibe-check create "Testing new feature"

# Create checkpoint with auto-generated timestamp
vibe-check create

# List all checkpoints  
vibe-check list

# Switch to specific checkpoint
vibe-check switch abc1234

# Finalize and push with custom message
vibe-check finalize "Implement user authentication"

# Finalize and push with auto-generated timestamp message
vibe-check finalize
```

### Available Commands

| Command | Description | Example |
|---------|-------------|---------|
| `vibe-check` | Launch interactive TUI | `vibe-check` |
| `vibe-check create [note]` | Create checkpoint with optional note | `vibe-check create "WIP: auth system"` |
| `vibe-check list` | Show all checkpoints with current marked | `vibe-check list` |
| `vibe-check switch <hash>` | Switch to specific checkpoint | `vibe-check switch abc1234` |
| `vibe-check finalize [message]` | Squash and push with optional message | `vibe-check finalize "Add login feature"` |
| `vibe-check --help` | Show all available commands | `vibe-check --help` |

### Auto-Generated Messages

When you don't provide custom notes or messages, vibe-check automatically generates them:

**Create checkpoint without note:**
```bash
vibe-check create
# Creates: "CHECKPOINT: 03/08/2025 14:30"
```

**Create checkpoint with note:**
```bash
vibe-check create "Testing auth"
# Creates: "CHECKPOINT: 03/08/2025 14:30 - Testing auth"
```

**Finalize without message:**
```bash
vibe-check finalize  
# Creates: "Update: 03/08/2025 14:30"
```

**Finalize with custom message:**
```bash
vibe-check finalize "Add user authentication"
# Creates: "Add user authentication"
```

### Simple Workflow (No Git Knowledge Required!)

1. **Make some changes** - Edit your code
2. **Create checkpoint** - Save a snapshot (arrow keys → Enter)
3. **Try different approach** - Make more changes, create another checkpoint
4. **Switch between versions** - Use "Change Checkpoint" to compare approaches
5. **Found the right solution?** - Use "Finalize and Push" to clean up and ship it

**That's it!** No complex Git commands, no fear of losing work, no messy commit history.

## 🆚 Before vs After

### Before (Traditional Git):
```bash
# Save work
git add .
git commit -m "WIP: trying approach 1"

# Try different approach  
git reset --hard HEAD~1
# Make changes
git add .
git commit -m "WIP: trying approach 2"

# Go back to first approach
git log --oneline  # find the commit hash
git checkout abc123

# Clean up commits
git rebase -i HEAD~5  # complex interactive editor
git push --force-with-lease origin main  # scary!
```

### After (vibe-check):
```bash
vibe-check
# → Arrow keys to navigate menu
# → Create Checkpoint (saves approach 1)
# → Create Checkpoint (saves approach 2) 
# → Change Checkpoint (switch between versions)
# → Finalize and Push (clean commit + push)
```

## 🎨 Screenshots

Beautiful terminal interface with:

<img width="369" height="248" alt="image" src="https://github.com/user-attachments/assets/b043d6c7-15f8-43e9-990e-ff3ad2624ab0" />

<img width="375" height="194" alt="image" src="https://github.com/user-attachments/assets/d5db9976-91ee-466a-b60c-321e7a964545" />

<img width="625" height="180" alt="image" src="https://github.com/user-attachments/assets/a25c2bba-1b5c-4926-974b-8b5dd3dff1db" />

<img width="513" height="179" alt="image" src="https://github.com/user-attachments/assets/6715a90a-1517-4aba-a31a-156a482bdd84" />

<img width="376" height="195" alt="image" src="https://github.com/user-attachments/assets/8b0d31bb-ba14-4918-a777-48129f145728" />

## 🛡️ Requirements

- Git repository (initialized)
- Go 1.21+ (for building from source)
- Terminal with color support

## 📖 Philosophy

**Git is powerful but intimidating.** Most developers spend more time fighting Git commands than actually coding. vibe-check fixes this by:

- **Hiding complexity** - You don't need to know `git rebase` or `git cherry-pick`
- **Preventing mistakes** - No more accidental `git reset --hard` disasters
- **Encouraging experimentation** - Try ideas fearlessly with easy checkpoints
- **Cleaning up automatically** - Messy development becomes clean commit history

Perfect for:
- 🤖 **AI-assisted coding** - Create checkpoints between AI suggestions
- 🧪 **Experimental development** - Try multiple approaches safely
- 🔄 **Refactoring** - Save checkpoints before big changes
- 📚 **Learning** - Experiment with new code without fear
- 👥 **Team projects** - Clean commit history for better collaboration

## 🤖 AI Assistant Instructions

Copy and paste this prompt into your AI coding session to teach your AI assistant about vibe-check:

<details>
<summary>📋 Click to expand AI Prompt (Copy-paste ready)</summary>

```text
VIBE-CHECK WORKFLOW INSTRUCTIONS:

You are working with a user who has vibe-check installed for Git checkpoint management. Follow these rules:

CHECKPOINT CREATION:
- After ANY code change that the user confirms works correctly, immediately create a checkpoint
- Use: vibe-check create "description of what was implemented/fixed"
- Be descriptive but concise in checkpoint descriptions
- Examples: "Add user authentication", "Fix database connection bug", "Implement email validation"

CHECKPOINT SWITCHING:
- BEFORE suggesting to switch to a previous checkpoint, ALWAYS ask: "Should I create a checkpoint of current progress first to avoid losing work?"
- Only suggest switching when comparing approaches or reverting to earlier working state
- Use: vibe-check switch <hash>
- Always explain WHY we're switching: "Let's switch back to checkpoint abc123 to compare the two authentication approaches"

FINALIZE AND PUSH:
- ONLY suggest vibe-check finalize when user explicitly wants to clean up and push to remote
- Always ask for confirmation: "Ready to finalize all checkpoints into a clean commit and push to remote?"
- Use custom message: vibe-check finalize "Meaningful commit message"

SAFETY RULES:
- NEVER switch checkpoints without creating a checkpoint first (risk of losing work)
- NEVER finalize without user's explicit permission
- Always suggest vibe-check list to show available checkpoints before switching
- If user seems lost, suggest vibe-check to open the interactive menu

WORKFLOW EXAMPLE:
1. User: "The login form works now"
   You: vibe-check create "Implement login form with validation"

2. User: "Let's try a different approach"
   You: "Should I create a checkpoint first? Then we can switch back if needed."
   User: "Yes"
   You: vibe-check create "Working login form - before trying alternative"

3. User: "I want to go back to the previous version"
   You: vibe-check list (show options), then vibe-check switch abc123

4. User: "This looks good, let's ship it"
   You: "Ready to finalize all checkpoints and push to remote?"
   User: "Yes"
   You: vibe-check finalize "Add complete login system with validation"
```

</details>

## 🤝 Contributing

**Vibe Check** is open source! We welcome contributions, bug reports, and feature requests.

- 🐛 **Found a bug?** [Open an issue](https://github.com/pr0d5h381/vibe-check/issues)
- 💡 **Have an idea?** [Start a discussion](https://github.com/pr0d5h381/vibe-check/discussions)
- 🛠️ **Want to contribute?** Fork the repo and submit a pull request!

### Development Setup

```bash
git clone https://github.com/pr0d5h381/vibe-check.git
cd vibe-check
go mod tidy
go build -o vibe-check .
```

## 📄 License

MIT License - feel free to use, modify, and distribute!

## 👨‍💻 Author

Created by **Adrian Górak** - [bomoge.pl](https://bomoge.pl)
