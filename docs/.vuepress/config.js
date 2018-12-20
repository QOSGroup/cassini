module.exports = {
    title: "Cassini Network",
    description: "Documentation for the Cassini Network.",
    dest: "./dist/docs",
    base: "/docs/",
    markdown: {
        lineNumbers: true
    },
    themeConfig: {
        lastUpdated: "Last Updated",
        nav: [{ text: "Back to Cassini", link: "https://cassini.network" }],
        sidebar: [
            {
                title: "Introduction",
                collapsable: false,
                children: [
                    "/",
                ]
            },
            {
                title: "Getting Started",
                collapsable: false,
                children: [
                    "/cassini.md",
                    "/quick_start.md",
                ]
            },
            {
                title: "Cassini Deployment",
                collapsable: false,
                children: [
                    "/cassini-deployment.md",
                ]
            }
            ,
            {
                title: "Etcd Config",
                collapsable: false,
                children: [
                    "/etcd_config.md",
                ]
            }

        ]
    }
}
