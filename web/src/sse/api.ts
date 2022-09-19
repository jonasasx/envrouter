import {Instance, InstancePod, Ref} from "../axios";

export interface SSEvent {
    itemType: "Ping" | "Instance" | "InstancePod" | "RefHead"
    item: Instance | InstancePod | Ref,
    event: "UPDATED" | "DELETED"
}