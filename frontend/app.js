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

const editError = document.getElementById("edit-error");
const createError = document.getElementById("create-error");

const API_URL = "";

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

function updateEditStatus() {
    editStatus.classList.toggle(
        "completed",
        editingTask.done
    );

    editTitle.classList.toggle(
        "completed",
        editingTask.done
    );
}

function getErrorMessage(errorCode) {
    switch (errorCode) {
        case "empty_title":
            return "Название задачи не может быть пустым.";

        case "invalid_task_dates":
            return "Дедлайн должен быть позже даты начала.";

        case "task_not_found":
            return "Задача не найдена.";

        default:
            return "Произошла неизвестная ошибка.";
    }
}

function openEditModal(task) {
    editingTask = task;
    editError.textContent = "";

    editTitle.value = task.title;
    editDescription.value = task.description || "";

    editScheduleDate.value = task.scheduleDate
        ? task.scheduleDate.slice(0, 10)
        : "";

    editDeadline.value = task.deadline
        ? task.deadline.slice(0, 10)
        : "";

    updateEditStatus();

    editModal.style.display = "flex";
}

async function editTaskTitle(task, titleElement) {
    const input = document.createElement("input");

    input.type = "text";
    input.value = task.title;

    input.className = "inline-title-input";

    titleElement.replaceWith(input);

    input.focus();
    input.select();

    const oldTitle = task.title;

    let finished = false;

    async function finish(save) {
        if (finished) {
            return;
        }

        finished = true;

        if (!save) {
            input.replaceWith(titleElement);
            return;
        }

        const newTitle = input.value.trim();

        if (!newTitle) {
            input.replaceWith(titleElement);
            return;
        }

        if (newTitle === oldTitle) {
            input.replaceWith(titleElement);
            return;
        }

        const response = await fetch(
            `${API_URL}/tasks/${task.id}`,
            {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    title: newTitle
                })
            }
        );

        if (!response.ok) {
            const data = await response.json();

            console.error(getErrorMessage(data.error));

            input.replaceWith(titleElement);

            return;
        }

        task.title = newTitle;
        titleElement.textContent = newTitle;

        input.replaceWith(titleElement);
    }

    input.addEventListener("keydown", async (event) => {
        if (event.key === "Enter") {
            await finish(true);
        }

        if (event.key === "Escape") {
            await finish(false);
        }
    });

    input.addEventListener("blur", async () => {
        await finish(true);
    });
}

async function loadTasks() {
    const url = currentFilter === "actual"
        ? `${API_URL}/tasks/actual`
        : `${API_URL}/tasks`;

    const response = await fetch(url);

    if (!response.ok) {
        const data = await response.json();

        console.error(
            "Не удалось загрузить задачи:",
            data.error
        );

        taskList.innerHTML =
            "<li class='empty-state'>Не удалось загрузить задачи.</li>";

        return;
    }

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
            <button class="done-button"></button>
        
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
        const taskTitleElement = li.querySelector(".task-info span");

        doneButton.addEventListener("click", async () => {
            const response = await fetch(
                `${API_URL}/tasks/${task.id}`,
                {
                    method: "PATCH",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({
                        done: !task.done
                    })
                }
            );

            if (!response.ok) {
                const data = await response.json();

                console.error(getErrorMessage(data.error));

                return;
            }

            await loadTasks();
        });

        deleteButton.addEventListener("click", async () => {
            const response = await fetch(
                `${API_URL}/tasks/${task.id}`,
                {
                    method: "DELETE"
                }
            );

            if (!response.ok) {
                const data = await response.json();

                console.error(getErrorMessage(data.error));

                return;
            }

            await loadTasks();
        });

        editButton.addEventListener("click", () => {
            editTaskTitle(task, taskTitleElement);
        });

        taskTitleElement.addEventListener("click", () => {
            openEditModal(task);
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
        `${API_URL}/tasks/${editingTask.id}`,
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
        const data = await response.json();

        editError.textContent = getErrorMessage(data.error);

        return;
    }

    editModal.style.display = "none";
    await loadTasks();
});

taskForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    createError.textContent = "";

    const title = taskTitle.value;

    const response = await fetch(`${API_URL}/tasks`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            title: title
        })
    });

    if (!response.ok) {
        const data = await response.json();

        createError.textContent = getErrorMessage(data.error);

        return;
    }

    taskTitle.value = "";
    await loadTasks();
});

deleteEdit.addEventListener("click", async () => {
    const response = await fetch(
        `${API_URL}/tasks/${editingTask.id}`,
        {
            method: "DELETE"
        }
    );

    if (!response.ok) {
        const data = await response.json();

        editError.textContent = getErrorMessage(data.error);

        return;
    }

    editModal.style.display = "none";
    await loadTasks();
});

editStatus.addEventListener("click", () => {
    editingTask.done = !editingTask.done;

    updateEditStatus();
});