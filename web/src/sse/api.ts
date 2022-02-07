import {Instance, InstancePod} from "../axios";

export interface SSEvent {
    itemType: "Ping" | "Instance" | "InstancePod"
    item: Instance | InstancePod,
    event: "UPDATED" | "DELETED"
}