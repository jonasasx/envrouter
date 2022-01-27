import {Card, Chip, Divider, Grid, List, ListItem} from "@mui/material";
import React from "react";
import {Theme} from "@mui/material/styles";
import {WithStyles, withStyles} from "@mui/styles";

interface PodProps extends WithStyles<typeof styles> {
}

const styles = (theme: Theme) => ({
    card: {
        marginBottom: "1em"
    }
})

export default withStyles(styles)(function PodComponent(props: PodProps) {
    const {classes} = props
    const pod = {
        pod: 'auth-deployment-dfd6b8cf4-fq6b7',
        name: 'auth-deployment',
        namespace: 'dev1',
        phase: 'Running',
        startedTime: '123',
        branch: 'master'
    }
    const rows: Array<{ key: string, value: any }> = [
        {key: "Pod name", value: pod.pod},
        {key: "App", value: pod.name},
        {key: "Env", value: pod.namespace},
        {key: "Shard", value: "s01"},
        {key: "Status", value: <Chip label={pod.phase} color={pod.phase === "Running" ? "success" : "warning"} variant="outlined" size="small"/>},
        {key: "Started", value: pod.startedTime},
        {key: "Branch", value: pod.branch || "master"},
        {key: "Commit", value: "c1873d26"},
        {key: "Author", value: "Ivan Volynkin"},
        {key: "Commit Message", value: "Init commit"},
    ]
    return (
        <Card className={classes.card}>
            <List>
                {
                    rows.map((row, index) => (
                        <React.Fragment key={row.key}>
                            {index === 0 || <Divider/>}
                            <ListItem>
                                <Grid container>
                                    <Grid item xs={4}>{row.key}:</Grid>
                                    <Grid item xs={8}>{row.value}</Grid>
                                </Grid>
                            </ListItem>
                        </React.Fragment>
                    ))
                }
            </List>
        </Card>
    )
})