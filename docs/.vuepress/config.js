module.exports = {
    title: "Cassini Network",
    description: "Documentation for the Cassini Network.",
    dest: "./dist/docs",
    base: "/cassini/",
    markdown: {
        lineNumbers: true
    },
    themeConfig: {
        lastUpdated: "Last Updated",
        nav: [{ text: "Back to Cassini", link: "http://docs.qoschain.info/cassini/" }],
        sidebarDepth:2,
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