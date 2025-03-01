async function saveNote() {
    const content = document.getElementById("newNote").value;
    const token = localStorage.getItem("access_token");

    const response = await fetch("http://localhost:8080/notes", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": token
        },
        body: JSON.stringify({ title: "Note", content })
    });

    if (response.ok) {
        alert("Note saved successfully!");
        loadNotes();
    } else {
        alert("Error saving note.");
    }
}

async function loadNotes() {
    const token = localStorage.getItem("access_token");

    const response = await fetch("http://localhost:8080/notes", {
        method: "GET",
        headers: { "Authorization": token }
    });

    const data = await response.json();
    console.log(data);

    if (data.length !== 0) {
        const notes = document.getElementById("notes");
        notes.innerHTML = "";

        if (response.ok) {
            data.forEach(note => {
                const noteElement = document.createElement("div");
                noteElement.innerHTML = `<h3>${note.title}</h3><p>${note.content}</p>`;
                notes.appendChild(noteElement);
            });
        } else {
            alert("Error loading notes.");
        }
    }
}

async function logout() {
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    window.location.href = "index.html";
}

document.addEventListener("DOMContentLoaded", () => {
    const token = localStorage.getItem("access_token");

    if (!token) {
        window.location.href = "index.html";
    }

    loadNotes();
});