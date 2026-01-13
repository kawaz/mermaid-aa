// Diagram types
export type DiagramType = 'flowchart' | 'sequence' | 'unknown';

// Flowchart types
export interface FlowchartNode {
  id: string;
  label: string;
  shape: 'rectangle' | 'rounded' | 'diamond' | 'circle' | 'stadium';
}

export interface FlowchartEdge {
  from: string;
  to: string;
  label?: string;
  style: 'solid' | 'dotted' | 'thick';
  arrow: 'arrow' | 'open' | 'none';
}

export interface FlowchartDiagram {
  type: 'flowchart';
  direction: 'TB' | 'TD' | 'BT' | 'LR' | 'RL';
  nodes: FlowchartNode[];
  edges: FlowchartEdge[];
}

// Sequence diagram types
export interface SequenceParticipant {
  id: string;
  label: string;
}

export interface SequenceMessage {
  from: string;
  to: string;
  label: string;
  style: 'solid' | 'dotted';
  arrow: 'arrow' | 'open';
}

export interface SequenceDiagram {
  type: 'sequence';
  participants: SequenceParticipant[];
  messages: SequenceMessage[];
}

export type Diagram = FlowchartDiagram | SequenceDiagram | { type: 'unknown' };
