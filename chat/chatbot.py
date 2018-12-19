from chatterbot import ChatBot
import sys, getopt


# Create a new instance of a ChatBot
bot = ChatBot(
    "Terminal",
    storage_adapter="chatterbot.storage.SQLStorageAdapter",
    logic_adapters=[
        "chatterbot.logic.BestMatch",
        "chatterbot.logic.MathematicalEvaluation",
        {
            'import_path': 'chatterbot.logic.SpecificResponseAdapter',
            'input_text': 'Help me',
            'output_text': 'Ok, plz contact us: 09096786789 '
        },
        {
            'import_path': 'chatterbot.logic.LowConfidenceAdapter',
            'threshold': 0.65,
            'default_response': 'I am sorry, but I do not understand.'
        }
    ],
    #input_adapter="chatterbot.input.TerminalAdapter",
    #output_adapter="chatterbot.output.TerminalAdapter",
    database="../database.db"
    #trainer="chatterbot.trainers.ChatterBotCorpusTrainer"
)


#bot.train('chatterbot.corpus.english')

def main(argv):
    try:
        opts, args = getopt.getopt(argv,"hq:",["question="])
    except getopt.GetoptError:
        print('test.py -i <inputfile> -o <outputfile>')
        sys.exit(2)
    for opt, arg in opts:
        if opt == '-h':
            print('cmd: python3 chatbot.py -q "hello"')
            sys.exit()
        elif opt in ("-q", "--question"):
            input = arg
            response = bot.get_response(input)
            print(response)
            sys.exit()


if __name__ == "__main__":
    main(sys.argv[1:])