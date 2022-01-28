import {InstancePod} from "../axios";

export interface PodEvent {
    item: InstancePod,
    event: "UPDATED" | "DELETED"
}