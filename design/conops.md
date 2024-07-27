# System Overview

Purpose of the proposed system or subsystem to which this document applies.
- General nature of the system
- Identify the:
    - Project sponsors
    - User agencies
    - Development organizations
    - Support agencies
    - Certifiers or certifying bodies
    - Operating centres or sites that will run the system
    - Document relevant to the present or proposed system

Graphical overview of the system is strongly recommended:
- Context diagram
- Top-level object diagram
- Type of diagram that depicts the system and its environment

Cited document:
- Project authorization
- Relevant technical document 
- Significant correspondence
- Documents concerning related projects
- Risk analysis reports
- Feasibility studies

# Referenced Documents

- List of document number, title, revision, and date of all referenced in this
  document.
- Identify source for all documents not available through normal channel. 

# Current System

System or situation (either automated or manual) as it currently exists:
- Watchtower
- CloudBOS DataDog
- Sentry
- Uptrace (staging)

If there is no current system on which to base changes:
- Describe situation that motivates development of the proposed system

Provide:
- Introduction to problem domain
- Enable readers to better understand the reasons for the desired changes and
  improvement.

# Background, Objectives and Scope

- Overview of the current system or situation (include if applicable):
    - Background
    - Mission
    - Objectives
    - Scope

When providing the background for the current system:
- Provides a brief summary of the motivation for the current system
    - Example:
        - Automation of certain tasks
        - Countering a certain threat situations
- Define the goal of the current system:
    - Strategies
    - Solutions
    - Tactics
    - Methods
    - Techniques used to accomplish the goal

Define the scope of the proposed system (brief summary only as it will discussed
in greater details in subsequent chapter):
- Mode of operation
- Class of users
- Interfaces to the operational environment 

# Operational Policies and Constraints

Operational policies and costraints that apply to the current system:
- Watchtower
- Sentry
- CloudBOS DataDog
- Hub Staging Uptrace

Operational Policies are predetermined management decisions regarding the
operations of the current system, normally in term of:
- General statements or understanding that guide decision making activities

Policies:
- Limit the decision making freedom but do allow for some discretion

Operational constraints include the following:
- A constraint on the hours of operation of the system, perhaps limited by
  access to secure terminals
- A constraint on the number of personnel available to operate the system
- A constraint on the computer hardware (must operate on computer X)
- A constraint on the operational facilities, such as office space

# Description of the current system

Description of the current system, including the following, as appropriate:
- Operational environment and its characteristics
- Major system components and the interconnection among those components
- Interfaces to external system / procedures.
- Capabilities, functions, and features of the current system.
- Charts and accompanying descriptions depicting:
    - inputs
    - outputs
    - data flows
    - control flows
    - manual and automated process
    - make it sufficient to understand the current system from user's point of
      view.
    - Cost of system operations
    - Operational risk factors
    - Performance characteristics such as speed, throughput, volume, frequency
    - Quality attributes, such as:
        - Availability
        - Correctness
        - Efficiency
        - Expendability
        - Flexibility
        - Interoperability
        - Maintainability
        - Portability
        - Reusability
        - Supportability
        - Usability
    - Provisions for safety, security, privacy, integrity, and continuity of
      operations in emergencies.

Since the purpose for this chapter is:
- To describe the current system and how it oeprates, it is appropriate to use
  any tools or techniques that serve this purpose.
- Important for the description to be simple enough and clear enough that all
  intended readers of the document can fully understand it.
- Write this document in term of user's terminology. Avoid terminology specific
  to computers (computer jargons).

Graphical tools should be used wherever possible.
- This document should be understandable by several different types of readers.
- Useful graphical tools:
    - Work breakdown strcuture
    - N2 charts
    - Sequence or activity charts
    - Functional flow block diagram
    - Structure charts
    - Allocation charts
    - Data flow diagrams
    - Object diagrams
    - Context diagrams
    - Storyboards
    - Entity relationship diagram

Description of the operational environment should: 
- identify, as applicable:
    - facilities
    - equipment
    - computing hardware
    - software
    - personnel
    - operating procedures used to operate the existing system
- the description should be detailed as necessary to give readers an
  understanding of the:
    - numbers
    - versions
    - capacity of the operational equipment being used
    - for example:
        - if the current system contain database:
            - capacity of the storage unit should be described, provided the
              information exerts an influence on the users' operational
              capabilities.
            - likewise, if the system uses communication links, the capacities
              of the communication links should be specified if the exert
              influence on factos such as user capabilities, response time or
              throughput

