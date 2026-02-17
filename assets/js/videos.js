function main() {
    const form = document.getElementById("upload-video");

    form.addEventListener("submit", async e => {
        e.preventDefault();

        const formData = new FormData(form);
        const url = "/api/users/videos/upload";

        try {
            const response = await fetch(url, {
                method: "POST",
                body: formData,
            });
            if (!response.ok) {
                const data = await response.text();
                throw new Error(data)
            }

            const data = await response.json();
            console.log(data);
        } catch (error) {
            console.log(error);
        }
    })
}

document.addEventListener("DOMContentLoaded", main);
