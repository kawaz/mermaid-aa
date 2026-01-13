import {
  Diagram,
  DiagramType,
  FlowchartDiagram,
  FlowchartNode,
  FlowchartEdge,
  SequenceDiagram,
  SequenceParticipant,
  SequenceMessage,
} from './types';

export function detectDiagramType(input: string): DiagramType {
  const trimmed = input.trim().toLowerCase();
  if (trimmed.startsWith('graph') || trimmed.startsWith('flowchart')) {
    return 'flowchart';
  }
  if (trimmed.startsWith('sequencediagram')) {
    return 'sequence';
  }
  return 'unknown';
}

export function parse(input: string): Diagram {
  const type = detectDiagramType(input);

  switch (type) {
    case 'flowchart':
      return parseFlowchart(input);
    case 'sequence':
      return parseSequence(input);
    default:
      return { type: 'unknown' };
  }
}

function parseFlowchart(input: string): FlowchartDiagram {
  // Split by newlines or semicolons
  const lines = input
    .split(/[\n;]/)
    .map(l => l.trim())
    .filter(l => l);
  const nodes: Map<string, FlowchartNode> = new Map();
  const edges: FlowchartEdge[] = [];

  // Parse direction from first line
  let direction: FlowchartDiagram['direction'] = 'TD';
  const firstLine = lines[0].toLowerCase();
  if (firstLine.includes(' lr')) direction = 'LR';
  else if (firstLine.includes(' rl')) direction = 'RL';
  else if (firstLine.includes(' bt')) direction = 'BT';
  else if (firstLine.includes(' tb') || firstLine.includes(' td')) direction = 'TD';

  // Parse nodes and edges
  for (let i = 1; i < lines.length; i++) {
    const line = lines[i];

    // Skip subgraph, end, style, class directives
    if (/^(subgraph|end|style|class|click|linkstyle)/i.test(line)) continue;

    // Parse edge patterns: A --> B, A --- B, A -.-> B, A ==> B, etc.
    const edgeMatch = line.match(
      /^(\w+)(\[[^\]]+\]|\([^)]+\)|\{[^}]+\}|>([^<]+)<\]|\(\([^)]+\)\)|\[\[[^\]]+\]\])?\s*(-->|---|-\.->|-.->|==>|-->\|[^|]+\||---\|[^|]+\||-\.->\|[^|]+\|)\s*(\w+)(\[[^\]]+\]|\([^)]+\)|\{[^}]+\}|>([^<]+)<\]|\(\([^)]+\)\)|\[\[[^\]]+\]\])?/
    );

    if (edgeMatch) {
      const fromId = edgeMatch[1];
      const fromShape = edgeMatch[2];
      const edgeStyle = edgeMatch[4];
      const toId = edgeMatch[5];
      const toShape = edgeMatch[6];

      // Add/update nodes
      if (!nodes.has(fromId)) {
        nodes.set(fromId, parseNodeShape(fromId, fromShape));
      } else if (fromShape && !nodes.get(fromId)!.label.includes(fromId)) {
        nodes.set(fromId, parseNodeShape(fromId, fromShape));
      }

      if (!nodes.has(toId)) {
        nodes.set(toId, parseNodeShape(toId, toShape));
      } else if (toShape && nodes.get(toId)!.label === toId) {
        nodes.set(toId, parseNodeShape(toId, toShape));
      }

      // Parse edge label
      let edgeLabel: string | undefined;
      const labelMatch = edgeStyle.match(/\|([^|]+)\|/);
      if (labelMatch) {
        edgeLabel = labelMatch[1];
      }

      // Determine edge style
      let style: FlowchartEdge['style'] = 'solid';
      let arrow: FlowchartEdge['arrow'] = 'arrow';

      if (edgeStyle.includes('-.') || edgeStyle.includes('-.-')) {
        style = 'dotted';
      } else if (edgeStyle.includes('==')) {
        style = 'thick';
      }

      if (edgeStyle.includes('---')) {
        arrow = 'none';
      }

      edges.push({ from: fromId, to: toId, label: edgeLabel, style, arrow });
      continue;
    }

    // Parse standalone node definition: A[Label] or A(Label) etc.
    const nodeMatch = line.match(/^(\w+)(\[[^\]]+\]|\([^)]+\)|\{[^}]+\}|\(\([^)]+\)\)|\[\[[^\]]+\]\])$/);
    if (nodeMatch) {
      const nodeId = nodeMatch[1];
      const nodeShape = nodeMatch[2];
      nodes.set(nodeId, parseNodeShape(nodeId, nodeShape));
    }
  }

  return {
    type: 'flowchart',
    direction,
    nodes: Array.from(nodes.values()),
    edges,
  };
}

function parseNodeShape(id: string, shapeStr?: string): FlowchartNode {
  if (!shapeStr) {
    return { id, label: id, shape: 'rectangle' };
  }

  let label = id;
  let shape: FlowchartNode['shape'] = 'rectangle';

  if (shapeStr.startsWith('((') && shapeStr.endsWith('))')) {
    // ((Circle))
    label = shapeStr.slice(2, -2);
    shape = 'circle';
  } else if (shapeStr.startsWith('([') && shapeStr.endsWith('])')) {
    // ([Stadium])
    label = shapeStr.slice(2, -2);
    shape = 'stadium';
  } else if (shapeStr.startsWith('[[') && shapeStr.endsWith(']]')) {
    // [[Subroutine]]
    label = shapeStr.slice(2, -2);
    shape = 'rectangle';
  } else if (shapeStr.startsWith('(') && shapeStr.endsWith(')')) {
    // (Rounded)
    label = shapeStr.slice(1, -1);
    shape = 'rounded';
  } else if (shapeStr.startsWith('{') && shapeStr.endsWith('}')) {
    // {Diamond}
    label = shapeStr.slice(1, -1);
    shape = 'diamond';
  } else if (shapeStr.startsWith('[') && shapeStr.endsWith(']')) {
    // [Rectangle]
    label = shapeStr.slice(1, -1);
    shape = 'rectangle';
  }

  return { id, label, shape };
}

function parseSequence(input: string): SequenceDiagram {
  // Split by newlines or semicolons
  const lines = input
    .split(/[\n;]/)
    .map(l => l.trim())
    .filter(l => l);
  const participants: Map<string, SequenceParticipant> = new Map();
  const messages: SequenceMessage[] = [];

  for (let i = 1; i < lines.length; i++) {
    const line = lines[i];

    // Skip notes, activations, loops, etc.
    if (/^(note|activate|deactivate|loop|end|alt|else|opt|par|and|rect|critical|break)/i.test(line)) {
      continue;
    }

    // Parse participant definition
    const participantMatch = line.match(/^participant\s+(\w+)(?:\s+as\s+(.+))?$/i);
    if (participantMatch) {
      const id = participantMatch[1];
      const label = participantMatch[2] || id;
      participants.set(id, { id, label });
      continue;
    }

    // Parse actor definition
    const actorMatch = line.match(/^actor\s+(\w+)(?:\s+as\s+(.+))?$/i);
    if (actorMatch) {
      const id = actorMatch[1];
      const label = actorMatch[2] || id;
      participants.set(id, { id, label });
      continue;
    }

    // Parse message: A->>B: Message or A-->>B: Message, A->B:, A-->B:
    const msgMatch = line.match(/^(\w+)\s*(--?>?>|--?>>?|-\.->?)\s*(\w+)\s*:\s*(.*)$/);
    if (msgMatch) {
      const from = msgMatch[1];
      const arrowStyle = msgMatch[2];
      const to = msgMatch[3];
      const label = msgMatch[4];

      // Auto-add participants
      if (!participants.has(from)) {
        participants.set(from, { id: from, label: from });
      }
      if (!participants.has(to)) {
        participants.set(to, { id: to, label: to });
      }

      const style: SequenceMessage['style'] = arrowStyle.includes('--') ? 'dotted' : 'solid';
      const arrow: SequenceMessage['arrow'] = arrowStyle.includes('>>') ? 'arrow' : 'open';

      messages.push({ from, to, label, style, arrow });
    }
  }

  return {
    type: 'sequence',
    participants: Array.from(participants.values()),
    messages,
  };
}
