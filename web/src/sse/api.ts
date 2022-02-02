import {Instance, InstancePod} from "../axios";

export interface SSEvent {
    itemType: "Instance" | "InstancePod"
    item: Instance | InstancePod,
    event: "UPDATED" | "DELETED"
}