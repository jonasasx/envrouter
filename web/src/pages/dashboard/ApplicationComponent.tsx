import {WithStyles, withStyles} from "@mui/styles";
import {Application, RefBinding} from "../../axios";
import {Theme} from "@mui/material/styles";
import {Grid, TextField} from "@mui/material";

interface ApplicationProps extends WithStyles<typeof styles> {
    application: Application
    refBinding: RefBinding | undefined
}

const styles = (theme: Theme) => ({
    row: {
        padding: ".6em 1em",
        "&:hover": {
            backgroundColor: theme.palette.action.hover
        }
    }
})

export default withStyles(styles)(function ApplicationComponent(props: ApplicationProps) {
    const {classes, application, refBinding} = props
    console.log(refBinding || '')
    return (
        <Grid className={classes.row} container>
            <Grid item xs={6} style={{display: "flex", alignItems: "flex-start"}}>
                <small>{application.name}</small>
            </Grid>
            <Grid item xs={6}>
                <TextField variant="standard" size="small" defaultValue={refBinding?.ref} />
            </Grid>
            <Grid item xs={12}>
            </Grid>
        </Grid>
    )
})