const taskList = document.getElementById("task-list");
const taskForm = document.getElementById("task-form");
const taskTitle = document.getElementById("task-title");

const editModal = document.getElementById("edit-modal");
const editTitle = document.getElementById("edit-title");
const editDescription = document.getElementById("edit-description");
const editScheduleDate = document.getElementById("edit-schedule-date");
const editDeadline = document.getElementById("edit-deadline");
const editStatus = document.getElementById("edit-status");

const cancelEdit = document.getElementById("cancel-edit");
const saveEdit = document.getElementById("save-edit");
const deleteEdit = document.getElementById("delete-edit");

const clearScheduleDate = document.getElementById("clear-schedule-date");
const clearDeadline = document.getElementById("clear-deadline");

let editingTask = null;

const allTasksButton = document.getElementById("all-tasks");
const actualTasksButton = document.getElementById("actual-tasks");

let currentFilter = "all";

cancelEdit.addEventListener("click", () => {
    editModal.style.display = "none";
});

clearScheduleDate.addEventListener("click", () => {
    editScheduleDate.value = "";
});

clearDeadline.addEventListener("click", () => {
    editDeadline.value = "";
});

function formatDate(date) {
    if (!date) {
        return "";
    }

    return new Date(date).toLocaleDateString("ru-RU");
}

function toApiDateTime(date) {
    if (!date) {
        return null;
    }

    return `${date}T00:00:00Z`;
}

function updateFilterButtons() {
    allTasksButton.classList.toggle(
        "active",
        currentFilter === "all"
    );

    actualTasksButton.classList.toggle(
        "active",
        currentFilter === "actual"
    );
}

async function loadTasks() {
    const url = currentFilter === "actual"
        ? "http://localhost:8080/tasks/actual"
        : "http://localhost:8080/tasks";

    const response = await fetch(url);
    const tasks = await response.json();

    taskList.innerHTML = "";

    if (tasks.length === 0) {
        taskList.innerHTML = "<li class='empty-state'>Пока нет задач ✨</li>";
        return;
    }

    tasks.forEach(task => {
        const li = document.createElement("li");
        if (task.done) {
            li.classList.add("completed");
        }

        li.innerHTML = `
            <button class="done-button">${task.done ? "✓" : ""}</button>
        
            <div class="task-info">
                <span>${task.title}</span>
                ${task.description ? `<small>${task.description}</small>` : ""}
                ${task.scheduleDate ? `<small>Начать: ${formatDate(task.scheduleDate)}</small>` : ""}
                ${task.deadline ? `<small>Дедлайн: ${formatDate(task.deadline)}</small>` : ""}
            </div>
        
            <button class="edit-button">✏️</button>
            <button class="delete-button">🗑</button>
        `;

        const doneButton = li.querySelector(".done-button");
        const deleteButton = li.querySelector(".delete-button");
        const editButton = li.querySelector(".edit-button");

        doneButton.addEventListener("click", async () => {
            await fetch(`http://localhost:8080/tasks/${task.id}`, {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    done: !task.done
                })
            });

            loadTasks();
        });

        deleteButton.addEventListener("click", async () => {
            await fetch(`http://localhost:8080/tasks/${task.id}`, {
                method: "DELETE"
            });

            loadTasks();
        });

        editButton.addEventListener("click", () => {
            editingTask = task;

            editTitle.value = task.title;
            editDescription.value = task.description || "";

            editScheduleDate.value = task.scheduleDate
                ? task.scheduleDate.slice(0, 10)
                : "";

            editDeadline.value = task.deadline
                ? task.deadline.slice(0, 10)
                : "";

            editStatus.textContent = task.done ? "✓" : "○";

            editModal.style.display = "flex";
        });

        taskList.appendChild(li);
    });
}

loadTasks();

allTasksButton.addEventListener("click", () => {
    currentFilter = "all";
    updateFilterButtons();
    loadTasks();
});

actualTasksButton.addEventListener("click", () => {
    currentFilter = "actual";
    updateFilterButtons();
    loadTasks();
});

saveEdit.addEventListener("click", async () => {
    if (!editingTask) {
        return;
    }

    const response = await fetch(
        `http://localhost:8080/tasks/${editingTask.id}`,
        {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                title: editTitle.value,
                description: editDescription.value,
                scheduleDate: toApiDateTime(editScheduleDate.value),
                deadline: toApiDateTime(editDeadline.value),
                done: editingTask.done
            })
        }
    );

    if (!response.ok) {
        const message = await response.text();
        console.error(
            "Ошибка сохранения:",
            response.status,
            message
        );
        return;
    }

    editModal.style.display = "none";
    await loadTasks();
});

taskForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    const title = taskTitle.value;

    await fetch("http://localhost:8080/tasks", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            title: title
        })
    });

    taskTitle.value = "";

    loadTasks();
});

deleteEdit.addEventListener("click", async () => {
    const response = await fetch(
        `http://localhost:8080/tasks/${editingTask.id}`,
        {
            method: "DELETE"
        }
    );

    if (!response.ok) {
        const message = await response.text();
        console.error("Ошибка удаления:", response.status, message);
        return;
    }

    editModal.style.display = "none";
    await loadTasks();
});

editStatus.addEventListener("click", () => {
    editingTask.done = !editingTask.done;

    editStatus.textContent = editingTask.done
        ? "✓"
        : "○";
});