require("cypress-failed-log");

export function TerminalLog(violations) {
    cy.task(
        'log',
        `${violations.length} accessibility violation${
            violations.length === 1 ? '' : 's'
        } ${violations.length === 1 ? 'was' : 'were'} detected`
    )

    const violationData = violations.map(
        ({ id, impact, description, nodes }) => ({
            id,
            impact,
            description,
            html: nodes[0].html,
            target: nodes[0].target[0]
        })
    )
    console.table(violationData)
    cy.task('table', violationData)
}