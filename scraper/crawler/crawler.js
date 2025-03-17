import puppeteer from "puppeteer";

const browser = await puppeteer.launch({ headless: false });
const page = await browser.newPage();

async function getBooks(url) {
    await page.goto(url);

    // Wait for the book table to load
    await page.waitForSelector(
        ".view-taxonomy-term-details-books .view-content table",
        { timeout: 5000 },
    )
        .catch(() => {
            return [];
        });

    // Extract books from the table using the correct selector
    const books = await page.$$eval(
        ".view-taxonomy-term-details-books .view-content table td.col-first a",
        (links) => {
            return links.map((link) => ({
                href: link.href,
                coverImage: link.querySelector("img")?.src || "",
                title: link.href.split("/").pop().replace(/-/g, " "), // Get the last part of the URL as title and clean it
            }));
        },
    ).catch(() => []);

    return books;
}

async function getSubcategories(url) {
    await page.goto(url);

    // Wait for subcategories to load with the correct selector
    await page.waitForSelector(".view-list-sub-tags .views-row a", {
        timeout: 5000,
    })
        .catch(() => {
            return [];
        });

    // Extract subcategories using the correct selector
    const subcategories = await page.$$eval(
        ".view-list-sub-tags .views-row a",
        (links) => {
            return links.map((link) => ({
                title: link.querySelector("h4").textContent.trim(),
                href: link.href,
            }));
        },
    ).catch(() => []);

    return subcategories;
}

async function main() {
    try {
        await page.goto("https://www.al-islam.org");

        // Click the categories dropdown
        await page.locator(
            ".active-trail.dropdown-toggle.active.has-submenu",
        ).click();

        // Wait for dropdown menu to be visible
        await page.waitForSelector(".dropdown-menu li a", { visible: true });

        // Get all links from the dropdown menu
        const mainCategories = await page.$$eval(
            ".dropdown-menu li a",
            (elements) => {
                return elements.map((el) => ({
                    href: el.href,
                    text: el.textContent.trim(),
                }));
            },
        );

        // Process each category
        for (const category of mainCategories) {
            console.log(`\n=== Main Category: ${category.text} ===`);
            console.log(`URL: ${category.href}`);

            const subcategories = await getSubcategories(category.href);

            if (subcategories.length > 0) {
                console.log("Subcategories:");
                for (const sub of subcategories) {
                    console.log(`\n  --- ${sub.title} ---`);
                    console.log(`  URL: ${sub.href}`);

                    const books = await getBooks(sub.href);

                    if (books.length > 0) {
                        console.log("  Books:");
                        books.forEach((book) => {
                            console.log(`    • ${book.title}`);
                            console.log(`      Link: ${book.href}`);
                            console.log(`      Cover: ${book.coverImage}`);
                        });
                    } else {
                        console.log("  No books found");
                    }
                }
            } else {
                console.log("No subcategories found");
            }
        }
    } catch (error) {
        console.error("Error:", error);
    } finally {
        await browser.close();
    }
}

main();
