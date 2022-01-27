import {Instance} from "../../axios";
import {withStyles, WithStyles} from "@mui/styles";
import {Theme} from "@mui/material/styles";
import {Chip, Dialog, DialogContent, DialogTitle} from "@mui/material";
import {useState} from "react";
import PodComponent from "./PodComponent";

interface InstanceProps extends WithStyles<typeof styles> {
    instance: Instance
}

const styles = (theme: Theme) => ({
    dialog: {
        backgroundColor: theme.palette.background.default,
    }
})

export default withStyles(styles)(function InstanceComponent(props: InstanceProps) {
    const {classes, instance} = props
    const [opened, setOpened] = useState(false)
    const color = "success"
    return (
        <span>
            <Chip label={instance.name} color={color} variant="outlined" size="small" onClick={() => setOpened(true)}/>
            <Dialog
                open={opened}
                onClose={() => setOpened(false)}
                fullWidth={true}
                scroll="body"
                aria-labelledby="alert-dialog-title"
                aria-describedby="alert-dialog-description"
            >
                <DialogTitle id="alert-dialog-title">
                    {`${instance.name}.${instance.environment}`}: 1
                </DialogTitle>
                <DialogContent dividers={true} className={classes.dialog}>
                    {
                        [0].map(pod => <PodComponent key={pod}/>)
                    }
                </DialogContent>
            </Dialog>
        </span>
    );
})